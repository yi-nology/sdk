package mongo

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"github.com/yi-nology/common/utils/xhash"
	"go.mongodb.org/mongo-driver/bson"
	option "go.mongodb.org/mongo-driver/mongo/options"
	"sort"
	"strings"
)

type MongoCli struct {
	*qmgo.QmgoClient
}

func (m Mongo) Init() (*MongoCli, error) {
	if m.Enable == false {
		return nil, nil
	}
	var opts []options.ClientOptions
	if m.IsZap {
		opts = m.GetClientOptions()
	}
	ctx := context.Background()
	client, err := qmgo.Open(ctx, &qmgo.Config{
		Uri:              m.Uri(),
		Coll:             m.Coll,
		Database:         m.Database,
		MinPoolSize:      &m.MinPoolSize,
		MaxPoolSize:      &m.MaxPoolSize,
		SocketTimeoutMS:  &m.SocketTimeoutMs,
		ConnectTimeoutMS: &m.ConnectTimeoutMs,
		Auth: &qmgo.Credential{
			Username:   m.Username,
			Password:   m.Password,
			AuthSource: m.AuthSource,
		},
	}, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "链接mongodb数据库失败!")
	}
	mc := MongoCli{
		QmgoClient: client,
	}
	if err = mc.Indexes(ctx); err != nil {
		return nil, err
	}
	return &mc, nil
}

type (
	Index struct {
		V    any      `bson:"v"`
		Ns   any      `bson:"ns"`
		Key  []bson.E `bson:"key"`
		Name string   `bson:"name"`
	}
)

func (m *MongoCli) Indexes(ctx context.Context) error {
	// 表名:索引列表 列: "表名": [][]string{{"index1", "index2"}}
	indexMap := map[string][][]string{}
	for collection, indexes := range indexMap {
		err := m.CreateIndexes(ctx, collection, indexes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MongoCli) CreateIndexes(ctx context.Context, name string, indexes [][]string) error {
	collection, err := m.Database.Collection(name).CloneCollection()
	if err != nil {
		return errors.Wrapf(err, "获取[%s]的表对象失败!", name)
	}
	list, err := collection.Indexes().List(ctx)
	if err != nil {
		return errors.Wrapf(err, "获取[%s]的索引对象失败!", name)
	}
	var entities []Index
	err = list.All(ctx, &entities)
	if err != nil {
		return errors.Wrapf(err, "获取[%s]的索引列表失败!", name)
	}
	length := len(indexes)
	indexMap1 := make(map[string][]string, length)
	for i := 0; i < length; i++ {
		sort.Strings(indexes[i]) // 对索引key进行排序, 在使用bson.M搜索时, bson会自动按照key的字母顺序进行排序
		length1 := len(indexes[i])
		keys := make([]string, 0, length1)
		for j := 0; j < length1; j++ {
			if indexes[i][i][0] == '-' {
				keys = append(keys, indexes[i][j], "-1")
				continue
			}
			keys = append(keys, indexes[i][j], "1")
		}
		key := strings.Join(keys, "_")
		_, o1 := indexMap1[key]
		if o1 {
			return errors.Errorf("索引[%s]重复!", key)
		}
		indexMap1[key] = indexes[i]
	}
	length = len(entities)
	indexMap2 := make(map[string]map[string]string, length)
	for i := 0; i < length; i++ {
		v1, o1 := indexMap2[entities[i].Name]
		if !o1 {
			keyLength := len(entities[i].Key)
			v1 = make(map[string]string, keyLength)
			for j := 0; j < keyLength; j++ {
				v2, o2 := v1[entities[i].Key[j].Key]
				if !o2 {
					v1 = make(map[string]string)
				}
				v2 = entities[i].Key[j].Key
				v1[entities[i].Key[j].Key] = v2
				indexMap2[entities[i].Name] = v1
			}
		}
	}
	for k1, v1 := range indexMap1 {
		_, o2 := indexMap2[k1]
		if o2 {
			continue
		} // 索引存在
		if len(fmt.Sprintf("%s.%s.$%s", collection.Name(), name, v1)) > 127 {
			err = m.Database.Collection(name).CreateOneIndex(ctx, options.IndexModel{
				Key:          v1,
				IndexOptions: option.Index().SetName(xhash.Md5Hex([]byte(k1))),
				// IndexOptions: option.Index().SetName(utils.MD5V([]byte(k1))).SetExpireAfterSeconds(86400), // SetExpireAfterSeconds(86400) 设置索引过期时间, 86400 = 1天
			})
			if err != nil {
				return errors.Wrapf(err, "创建索引[%s]失败!", k1)
			}
			return nil
		}
		err = m.Database.Collection(name).CreateOneIndex(ctx, options.IndexModel{
			Key:          v1,
			IndexOptions: option.Index().SetExpireAfterSeconds(86400),
			// IndexOptions: option.Index().SetName(utils.MD5V([]byte(k1))).SetExpireAfterSeconds(86400), // SetExpireAfterSeconds(86400) 设置索引过期时间(秒), 86400 = 1天
		})
		if err != nil {
			return errors.Wrapf(err, "创建索引[%s]失败!", k1)
		}
	}
	return nil
}
