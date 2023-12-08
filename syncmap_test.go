package syncmap

import (
	"reflect"
	"sync"
	"testing"
)

func Test_syncmapImpl_Get(t *testing.T) {
	type fields struct {
		syncmap *sync.Map
	}
	type args struct {
		key any
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			syncmap := &sync.Map{}
			syncmap.Store(1, "test")
			return testCase{
				name: "number:string",
				fields: fields{
					syncmap: syncmap,
				},
				args: args{
					key: 1,
				},
				want:    "test",
				wantErr: false,
			}
		}(),
		func() testCase {
			syncmap := &sync.Map{}
			syncmap.Store("abc", 'x')
			return testCase{
				name: "string:rune",
				fields: fields{
					syncmap: syncmap,
				},
				args: args{
					key: "abc",
				},
				want:    'x',
				wantErr: false,
			}
		}(),
		func() testCase {
			syncmap := &sync.Map{}
			syncmap.Store('a', "🥺") // U+1F97A
			return testCase{
				name: "rune:emoji",
				fields: fields{
					syncmap: syncmap,
				},
				args: args{
					key: 'a',
				},
				want:    "🥺",
				wantErr: false,
			}
		}(),
		func() testCase {
			syncmap := &sync.Map{}
			syncmap.Store("🫠", 100) // U+1FAE0
			return testCase{
				name: "emoji:number",
				fields: fields{
					syncmap: syncmap,
				},
				args: args{
					key: "🫠",
				},
				want:    100,
				wantErr: false,
			}
		}(),
		func() testCase {
			syncmap := &sync.Map{}
			return testCase{
				name: "not found",
				fields: fields{
					syncmap: syncmap,
				},
				args: args{
					key: 1,
				},
				wantErr: true,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &syncmapImpl{
				syncmap: tt.fields.syncmap,
			}
			got, err := impl.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_syncmapImpl_GetMulti(t *testing.T) {
	type fields struct {
		syncmap *sync.Map
	}
	type args struct {
		keys []any
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		want    map[any]any
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			syncmap := &sync.Map{}
			syncmap.Store(1, "test")
			syncmap.Store("abc", 'x')
			syncmap.Store('a', "🥺")
			syncmap.Store("🫠", 100)
			return testCase{
				name: "normal",
				fields: fields{
					syncmap: syncmap,
				},
				args: args{
					keys: []any{1, "abc", 'a', "🫠"},
				},
				want: map[any]any{
					1:     "test",
					"abc": 'x',
					'a':   "🥺",
					"🫠":   100,
				},
				wantErr: false,
			}
		}(),
		func() testCase {
			syncmap := &sync.Map{}
			syncmap.Store(1, "test")
			syncmap.Store('a', "🥺")
			syncmap.Store("🫠", 100)
			return testCase{
				name: "including not found",
				fields: fields{
					syncmap: syncmap,
				},
				args: args{
					keys: []any{1, "abc", 'a', "🫠"},
				},
				want: map[any]any{
					1:   "test",
					'a': "🥺",
					"🫠": 100,
				},
				wantErr: false,
			}
		}(),
		func() testCase {
			syncmap := &sync.Map{}
			return testCase{
				name: "blank",
				fields: fields{
					syncmap: syncmap,
				},
				args: args{
					keys: []any{1, "abc", 'a', "🫠"},
				},
				want:    map[any]any{},
				wantErr: false,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &syncmapImpl{
				syncmap: tt.fields.syncmap,
			}
			got, err := impl.GetMulti(tt.args.keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMulti() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMulti() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_syncmapImpl_Create(t *testing.T) {
	type fields struct {
		syncmap *sync.Map
	}
	type args struct {
		key   any
		value any
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			syncmap := &sync.Map{}
			return testCase{
				name: "normal",
				fields: fields{
					syncmap: syncmap,
				},
				args: args{
					key:   1,
					value: "test",
				},
				wantErr: false,
			}
		}(),
		func() testCase {
			syncmap := &sync.Map{}
			return testCase{
				name: "zero value",
				fields: fields{
					syncmap: syncmap,
				},
				args: args{
					key:   0,
					value: "",
				},
				wantErr: false,
			}
		}(),
		func() testCase {
			syncmap := &sync.Map{}
			return testCase{
				name: "key is nil",
				fields: fields{
					syncmap: syncmap,
				},
				args: args{
					key:   nil,
					value: "",
				},
				wantErr: true,
			}
		}(),
		func() testCase {
			syncmap := &sync.Map{}
			return testCase{
				name: "value is nil",
				fields: fields{
					syncmap: syncmap,
				},
				args: args{
					key:   0,
					value: nil,
				},
				wantErr: true,
			}
		}(),
		func() testCase {
			syncmap := &sync.Map{}
			syncmap.Store(1, "test")
			return testCase{
				name: "already exists",
				fields: fields{
					syncmap: syncmap,
				},
				args: args{
					key:   1,
					value: "test",
				},
				wantErr: true,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &syncmapImpl{
				syncmap: tt.fields.syncmap,
			}
			if err := impl.Create(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			value, ok := tt.fields.syncmap.Load(tt.args.key)
			if !ok {
				t.Errorf("Create() failed to store args = %v", tt.args)
				return
			}
			if value != tt.args.value {
				t.Errorf("Create() failed to store args = %v, got = %v", tt.args, value)
			}
		})
	}
}

func Test_syncmapImpl_Update(t *testing.T) {
	type fields struct {
		syncmap *sync.Map
	}
	type args struct {
		key   any
		value any
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			syncmap := &sync.Map{}
			syncmap.Store(1, "test")
			return testCase{
				name: "normal",
				fields: fields{
					syncmap: syncmap,
				},
				args: args{
					key:   1,
					value: "test2",
				},
				wantErr: false,
			}
		}(),
		func() testCase {
			syncmap := &sync.Map{}
			return testCase{
				name: "not found caused no error",
				fields: fields{
					syncmap: syncmap,
				},
				args: args{
					key:   1,
					value: "test2",
				},
				wantErr: false,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &syncmapImpl{
				syncmap: tt.fields.syncmap,
			}
			if err := impl.Update(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			value, ok := tt.fields.syncmap.Load(tt.args.key)
			if !ok {
				t.Errorf("Update() failed to store args = %v", tt.args)
				return
			}
			if value != tt.args.value {
				t.Errorf("Update() failed to store args = %v, got = %v", tt.args, value)
			}
		})
	}
}

func Test_syncmapImpl_Delete(t *testing.T) {
	type fields struct {
		syncmap *sync.Map
	}
	type args struct {
		key any
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			syncmap := &sync.Map{}
			syncmap.Store(1, "test")
			return testCase{
				name: "normal",
				fields: fields{
					syncmap: syncmap,
				},
				args: args{
					key: 1,
				},
				wantErr: false,
			}
		}(),
		func() testCase {
			syncmap := &sync.Map{}
			return testCase{
				name: "not found is not error",
				fields: fields{
					syncmap: syncmap,
				},
				args: args{
					key: 1,
				},
				wantErr: false,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &syncmapImpl{
				syncmap: tt.fields.syncmap,
			}
			if err := impl.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
