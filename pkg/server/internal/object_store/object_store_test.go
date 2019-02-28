package object_store

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var objstore *ObjectStore

func prepareObjectStore() {

	var (
		appid  = os.Getenv("AWSAPPID")
		secret = os.Getenv("AWSSECRET")
		region = os.Getenv("AWSREGION")
		bucket = os.Getenv("AWSBUCKET")
		root   = os.Getenv("AWSROOT")
	)

	objstore = NewObjectStore(appid,
		secret,
		Region(region),
		StoreBucket(bucket),
		Root(root),
	)

}

func TestMain(t *testing.T) {
	prepareObjectStore()
}

func TestListBuckets(t *testing.T) {
	prepareObjectStore()

	buckets, err := objstore.ListBuckets()
	if err != nil {
		t.Fatalf("ListBuckets error %s", err)
	}

	if len(buckets) == 0 {
		t.Fatal("Don't has buckets")
	}

	t.Log("Buckets ---")
	for _, buck := range buckets {
		t.Logf("\t%s", buck)
	}

}

func TestListObjects(t *testing.T) {
	prepareObjectStore()
	objects, err := objstore.ListObjects("")
	if err != nil {
		t.Fatalf("ListObjects error %s", err)
	}

	if len(objects) == 0 {
		t.Fatal("Don't has objects")
	}

	t.Log("Objects ---")
	for _, obj := range objects {
		t.Logf("\t%s", obj)
	}

	objects, err = objstore.ListObjects("datasets")
	if err != nil {
		t.Fatalf("ListObjects error %s", err)
	}

	if len(objects) == 0 {
		t.Fatal("Don't has objects")
	}

	t.Log("Objects ---")
	for _, obj := range objects {
		t.Logf("\t%s", obj)
	}
}

func TestReadWriteObject(t *testing.T) {
	prepareObjectStore()
	filename := "test/TEST.txt"
	size, version, err := objstore.PutObject(filename, []byte("hello world"))
	if err != nil {
		t.Fatalf("PutObject error %s", err)
	}

	t.Logf("size %d, version %s", size, version)

	body, err := objstore.GetObject(filename, 0, 0)
	if err != nil {
		t.Fatalf("GetObject error %s", err)
	}

	content, err := ioutil.ReadAll(body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll error %s", err)
	}

	t.Logf("Contents: %s", content)
}

func TestDeleteObject(t *testing.T) {
	prepareObjectStore()
	filename := "test/deleted.txt"
	size, version, err := objstore.PutObject(filename, []byte("must be deleted"))
	if err != nil {
		t.Fatalf("PutObject error %s", err)
	}

	t.Logf("size %d, version %s", size, version)

	mark, err := objstore.DeleteObject(filename)
	if err != nil {
		t.Fatalf("DeleteObject error %s", err)
	}

	t.Logf("Delete Marker %v", mark)
}

func TestCopyObject(t *testing.T) {
	prepareObjectStore()
	src := "test/1.txt"
	size, version, err := objstore.PutObject(src, []byte("file content"))
	if err != nil {
		t.Fatalf("PutObject error %s", err)
	}

	t.Logf("size %d, version %s", size, version)

	err = objstore.CopyObject(src, "test/2.txt")
	if err != nil {
		t.Fatalf("CopyObject error %s", err)
	}

	t.Logf("Copy Success")
}

func TestGetObjectPresignURL(t *testing.T) {
	prepareObjectStore()
	filename := "test/get_sign.txt"
	size, version, err := objstore.PutObject(filename, []byte("need download content"))
	if err != nil {
		t.Fatalf("PutObject error %s", err)
	}

	t.Logf("size %d, version %s", size, version)

	url, signedHeaders, err := objstore.GetObjectPresignURL(filename)
	if err != nil {
		t.Fatalf("GetObjectPresignURL error %s", err)
	}

	t.Logf("URL: %s, %v", url, signedHeaders)
}

func TestPutObjectPresignURL(t *testing.T) {
	prepareObjectStore()
	filename := "test/put_sign.txt"

	url, signedHeaders, err := objstore.PutObjectPresignURL(filename)
	if err != nil {
		t.Fatalf("PutObjectPresignURL error %s", err)
	}

	t.Logf("URL: %s, %v", url, signedHeaders)
}

func TestMoveObject(t *testing.T) {
	prepareObjectStore()
	src := "test/3.txt"
	size, version, err := objstore.PutObject(src, []byte("file content"))
	if err != nil {
		t.Fatalf("PutObject error %s", err)
	}

	t.Logf("size %d, version %s", size, version)

	err = objstore.MoveObject(src, "test/4.txt")
	if err != nil {
		t.Fatalf("MoveObject error %s", err)
	}

	t.Logf("MoveObject Success")
}

func TestExistsObject(t *testing.T) {
	prepareObjectStore()
	filename := "test/5.txt"
	size, version, err := objstore.PutObject(filename, []byte("file content"))
	if err != nil {
		t.Fatalf("PutObject error %s", err)
	}

	t.Logf("size %d, version %s", size, version)

	ok, err := objstore.ExistsObject(filename)
	if err != nil {
		t.Fatalf("ExistsObject error %s", err)
	}

	if ok {
		t.Logf("ExistsObject Success")
	}

	_, err = objstore.DeleteObject(filename)
	if err != nil {
		t.Fatalf("DeleteObject error %s", err)
	}

	ok, err = objstore.NotExistsObject(filename)
	if err != nil {
		t.Fatalf("NotExistsObject error %s", err)
	}

	if ok {
		t.Logf("NotExistsObject Success")
	}
}

func TestPutObjectManager(t *testing.T) {
	prepareObjectStore()
	rd := strings.NewReader("file content")
	src := "test/upload.txt"
	version, err := objstore.PutObjectManager(src, rd)
	if err != nil {
		t.Fatalf("PutObject error %s", err)
	}

	t.Logf("PutObjectManager success version %s", version)

}

func TestLargePutObjectManager(t *testing.T) {
	prepareObjectStore()
	var buf = new(bytes.Buffer)
	buf.Grow(1024 * 1204)
	go func() {
		for i := 0; i < 100; i++ {
			buf.Write(make([]byte, 1000))
		}
	}()

	src := "test/large.txt"
	version, err := objstore.PutObjectManager(src, buf)
	if err != nil {
		t.Fatalf("PutObjectManager error %s", err)
	}

	t.Logf("PutObjectManager success version %s", version)
}

func TestGetObjectManager(t *testing.T) {
	prepareObjectStore()
	var buf = new(WriteAtBuffer)

	src := "test/large.txt"

	size, err := objstore.GetObjectManager(src, buf)
	if err != nil {
		t.Fatalf("GetObjectManager error %s", err)
	}

	t.Logf("GetObjectManager success size %d", size)
}
