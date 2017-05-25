package ds

import (
	"context"
)

func beforeSave(src interface{}) {
	x := src.(KeyGetSetter)

	// stamp model
	if x, ok := src.(Stampable); ok {
		x.Stamp()
	}

	// create new key
	if x.GetKey() == nil {
		if x, ok := src.(KeyNewer); ok {
			x.NewKey()
		}
	}
}

// SaveModel saves model to datastore
// if key was not set in model, will call NewKey
func (client *Client) SaveModel(ctx context.Context, src interface{}) error {
	beforeSave(src)
	err := client.PutModel(ctx, src)
	if err != nil {
		return err
	}
	return nil
}

// SaveModels saves models to datastore
// see more in SaveModel
func (client *Client) SaveModels(ctx context.Context, src interface{}) error {
	xs := valueOf(src)
	for i := 0; i < xs.Len(); i++ {
		x := xs.Index(i).Interface()
		beforeSave(x)
	}
	err := client.PutModels(ctx, src)
	if err != nil {
		return err
	}
	return nil
}
