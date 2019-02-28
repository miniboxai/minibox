package object_store

import "fmt"

func (b *Bucket) String() string {
	return fmt.Sprintf("%s - %s", b.CreatedAt, b.Name)
}

func (obj *Object) String() string {
	return fmt.Sprintf("%s %8d - %s", obj.LastModified, obj.Size, obj.Key)
}
