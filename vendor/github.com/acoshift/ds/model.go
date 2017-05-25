package ds

import (
	"strconv"

	"cloud.google.com/go/datastore"
)

// KeyGetter interface
type KeyGetter interface {
	GetKey() *datastore.Key
}

// KeySetter interface
type KeySetter interface {
	SetKey(*datastore.Key)
}

// KeyGetSetter interface
type KeyGetSetter interface {
	KeyGetter
	KeySetter
}

// KeyNewer interface
type KeyNewer interface {
	NewKey()
}

// Model is the base model which id is int64
type Model struct {
	Key *datastore.Key `datastore:"__key__"`
}

// GetKey returns key from model
func (x *Model) GetKey() *datastore.Key {
	if x == nil {
		return nil
	}
	return x.Key
}

// SetKey sets model key to given key
func (x *Model) SetKey(key *datastore.Key) {
	if x == nil {
		return
	}
	x.Key = key
}

// SetID sets id key to model
func (x *Model) SetID(kind string, id int64) {
	x.SetKey(datastore.IDKey(kind, id, nil))
}

// ID returns id
func (x *Model) ID() int64 {
	if x == nil || x.Key == nil {
		return 0
	}
	return x.Key.ID
}

// StringID return id in string format
func (x *Model) StringID() string {
	if x == nil || x.Key == nil {
		return ""
	}
	return strconv.FormatInt(x.Key.ID, 10)
}

// NewIncomplateKey sets an incomplete key to model
func (x *Model) NewIncomplateKey(kind string, parent *datastore.Key) {
	x.SetKey(datastore.IncompleteKey(kind, parent))
}

// NewKey sets incomplete key to model
// func (x *Model) NewKey() {
// 	x.SetKey(datastore.IncompleteKey(kind, nil))
// }

// StringIDModel is the base model which id is string
// but can use both id key and name key
type StringIDModel struct {
	Key *datastore.Key `datastore:"__key__"`
}

// GetKey returns key from model
func (x *StringIDModel) GetKey() *datastore.Key {
	if x == nil {
		return nil
	}
	return x.Key
}

// SetKey sets model key to given key
// if key is not name key, it will use id key
func (x *StringIDModel) SetKey(key *datastore.Key) {
	if x == nil {
		return
	}
	x.Key = key
}

// SetID sets id to model
func (x *StringIDModel) SetID(kind string, id int64) {
	if id == 0 {
		// invalid key
		// if set id key to 0, datastore server will throw error, which we can not handle
		return
	}
	x.SetKey(datastore.IDKey(kind, id, nil))
}

// SetStringID sets string id to model
func (x *StringIDModel) SetStringID(kind string, id string) {
	x.SetID(kind, parseID(id))
}

// SetNameID sets name id to model
func (x *StringIDModel) SetNameID(kind string, name string) {
	x.SetKey(datastore.NameKey(kind, name, nil))
}

// NewIncomplateKey sets an incomplete key to model
func (x *StringIDModel) NewIncomplateKey(kind string, parent *datastore.Key) {
	x.SetKey(datastore.IncompleteKey(kind, parent))
}

// NewKey sets incomplete key to model
// func (x *StringIDModel) NewKey(kind string) {
// 	x.SetKey(datastore.IncompleteKey(kind, nil))
// }

// ID return id
func (x *StringIDModel) ID() string {
	if x == nil || x.Key == nil {
		return ""
	}
	if x.Key.Name != "" {
		return x.Key.Name
	}
	if x.Key.ID != 0 {
		return strconv.FormatInt(x.Key.ID, 10)
	}
	return ""
}

// SetKey sets key to model
func SetKey(key *datastore.Key, dst interface{}) {
	if dst == nil || key == nil {
		return
	}
	if x, ok := dst.(KeySetter); ok {
		x.SetKey(key)
	}
}

// SetKeys sets keys to models
func SetKeys(keys []*datastore.Key, dst interface{}) {
	if dst == nil || len(keys) == 0 {
		return
	}
	xs := valueOf(dst)
	for i := 0; i < xs.Len(); i++ {
		x := xs.Index(i)
		if x.IsNil() {
			continue
		}
		if x, ok := x.Interface().(KeySetter); ok {
			x.SetKey(keys[i])
		}
	}
}

// SetCommitKey sets commit pending key to model
func SetCommitKey(commit *datastore.Commit, pendingKey *datastore.PendingKey, dst interface{}) {
	if dst == nil {
		return
	}
	if x, ok := dst.(KeySetter); ok {
		x.SetKey(commit.Key(pendingKey))
	}
}

// SetCommitKeys sets commit pending keys to models
func SetCommitKeys(commit *datastore.Commit, pendingKeys []*datastore.PendingKey, dst interface{}) {
	xs := valueOf(dst)
	for i := 0; i < xs.Len(); i++ {
		x := xs.Index(i)
		if x.IsNil() {
			continue
		}
		if x, ok := x.Interface().(KeySetter); ok {
			x.SetKey(commit.Key(pendingKeys[i]))
		}
	}
}

// SetID sets id to model
func SetID(kind string, id int64, dst interface{}) {
	SetKey(datastore.IDKey(kind, id, nil), dst)
}

// SetIDs sets ids to models
func SetIDs(kind string, ids []int64, dst interface{}) {
	keys := make([]*datastore.Key, len(ids))
	for i := range ids {
		keys[i] = datastore.IDKey(kind, ids[i], nil)
	}
	SetKeys(keys, dst)
}

// SetStringID sets string id to model
func SetStringID(kind string, id string, dst interface{}) {
	tid := parseID(id)
	if tid == 0 {
		return
	}
	SetKey(datastore.IDKey(kind, tid, nil), dst)
}

// SetStringIDs sets string id to models
func SetStringIDs(kind string, ids []string, dst interface{}) {
	keys := make([]*datastore.Key, len(ids))
	for i := range ids {
		keys[i] = datastore.IDKey(kind, parseID(ids[i]), nil)
	}
	SetKeys(keys, dst)
}

// SetNameIDs sets name id to models
func SetNameIDs(kind string, names []string, dst interface{}) {
	keys := make([]*datastore.Key, len(names))
	for i := range names {
		keys[i] = datastore.NameKey(kind, names[i], nil)
	}
	SetKeys(keys, dst)
}
