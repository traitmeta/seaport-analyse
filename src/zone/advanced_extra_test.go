package zone

import (
	"testing"
)

func TestAdvancedExtra_toPaddedBytes(t *testing.T) {
	type fields struct {
		contract string
		chainId  string
	}
	type args struct {
		value    uint64
		numBytes int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "test expiretime",
			fields: fields{
				contract: "",
				chainId:  "",
			},
			args: args{
				value:    10000000,
				numBytes: 8,
			},
			want: "0000000000989680",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &AdvancedExtraBuilder{
				contract: tt.fields.contract,
				chainId:  tt.fields.chainId,
			}
			if got := b.toPaddedBytes(tt.args.value, tt.args.numBytes); got != tt.want {
				t.Errorf("AdvancedExtraBiz.toPaddedBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdvancedExtra_GetZoneContext(t *testing.T) {
	type fields struct {
		contract  string
		chainId   string
		fulfiller string
	}
	type args struct {
		considerFirstItemIdentifier string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "test 1",
			fields: fields{
				contract:  "",
				chainId:   "",
				fulfiller: "",
			},
			args: args{
				considerFirstItemIdentifier: "1",
			},
			want: "0x000000000000000000000000000000000000000000000000000000000000000001",
		},
		{
			name: "test 10",
			fields: fields{
				contract:  "",
				chainId:   "",
				fulfiller: "",
			},
			args: args{
				considerFirstItemIdentifier: "10",
			},
			want: "0x00000000000000000000000000000000000000000000000000000000000000000a",
		},
		{
			name: "test 0x1",
			fields: fields{
				contract:  "",
				chainId:   "",
				fulfiller: "",
			},
			args: args{
				considerFirstItemIdentifier: "0x1",
			},
			want: "0x000000000000000000000000000000000000000000000000000000000000000001",
		},
		{
			name: "test 0x1b",
			fields: fields{
				contract:  "",
				chainId:   "",
				fulfiller: "",
			},
			args: args{
				considerFirstItemIdentifier: "0x1b",
			},
			want: "0x00000000000000000000000000000000000000000000000000000000000000001b",
		},
		{
			name: "test 0x8b02e2c4613ed1cd1d15021beedc8dc421f65e22feb5ffe2468d386a3914cae1",
			fields: fields{
				contract:  "",
				chainId:   "",
				fulfiller: "",
			},
			args: args{
				considerFirstItemIdentifier: "0x8b02e2c4613ed1cd1d15021beedc8dc421f65e22feb5ffe2468d386a3914cae1",
			},
			want: "0x008b02e2c4613ed1cd1d15021beedc8dc421f65e22feb5ffe2468d386a3914cae1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &AdvancedExtraBuilder{
				contract:  tt.fields.contract,
				chainId:   tt.fields.chainId,
				fulfiller: tt.fields.fulfiller,
			}
			if got, success := b.BuildZoneContext(tt.args.considerFirstItemIdentifier); !success || got != tt.want {
				t.Errorf("AdvancedExtraBiz.GetZoneContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdvancedExtra_getExtraData(t *testing.T) {
	type fields struct {
		contract  string
		chainId   string
		fulfiller string
	}
	type args struct {
		priv       string
		orderHash  string
		context    string
		expiration int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test get extra data",
			fields: fields{
				contract:  "0x23cafcac35dee6f487493f7eea284d8689b8c179",
				chainId:   "0x94a9059e",
				fulfiller: "0x7a08b9d8d847c91f81eb7d6726fd86fc577183bd",
			},
			args: args{
				priv:       "7e5bfb82febc4c2c8529167104271ceec190eafdca277314912eaabdb67c6e5f",
				orderHash:  "0x8b02e2c4613ed1cd1d15021beedc8dc421f65e22feb5ffe2468d386a3914cae1",
				context:    "0x00000000000000000000000000000000000000000000000000000000000000000c",
				expiration: 1755121549,
			},
			want:    "0x007a08b9d8d847c91f81eb7d6726fd86fc577183bd00000000689d078dd467dc66314693d92dea1609924949439d03bc962821cabd5a4bbd9d4921d8f4d3e3f8b7f5ff9e989dffa346165b92e9e55573225636c08d4045c6a81fc1e59400000000000000000000000000000000000000000000000000000000000000000c",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &AdvancedExtraBuilder{
				contract:  tt.fields.contract,
				chainId:   tt.fields.chainId,
				fulfiller: tt.fields.fulfiller,
			}
			got, err := b.buildExtraData(tt.args.priv, tt.args.orderHash, tt.args.context, tt.args.expiration)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdvancedExtraBiz.getExtraData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AdvancedExtraBiz.getExtraData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdvancedExtra_convertSignatureToEIP2098(t *testing.T) {
	type fields struct {
		contract string
		chainId  string
	}
	type args struct {
		signature string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "test for yParity = 0",
			fields: fields{
				contract: "",
				chainId:  "",
			},
			args: args{
				signature: "0x68a020a209d3d56c46f38cc50a33f704f4a9a10a59377f8dd762ac66910e9b907e865ad05c4035ab5792787d4a0297a43617ae897930a6fe4d822b8faea520641b",
			},
			want: "0x68a020a209d3d56c46f38cc50a33f704f4a9a10a59377f8dd762ac66910e9b907e865ad05c4035ab5792787d4a0297a43617ae897930a6fe4d822b8faea52064",
		},
		{
			name: "test for yParity = 1",
			fields: fields{
				contract: "",
				chainId:  "",
			},
			args: args{
				signature: "0x9328da16089fcba9bececa81663203989f2df5fe1faa6291a45381c81bd17f76139c6d6b623b42da56557e5e734a43dc83345ddfadec52cbe24d0cc64f5507931c",
			},
			want: "0x9328da16089fcba9bececa81663203989f2df5fe1faa6291a45381c81bd17f76939c6d6b623b42da56557e5e734a43dc83345ddfadec52cbe24d0cc64f550793",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &AdvancedExtraBuilder{
				contract: tt.fields.contract,
				chainId:  tt.fields.chainId,
			}
			if got := b.convertSignatureToEIP2098(tt.args.signature); got != tt.want {
				t.Errorf("AdvancedExtraBiz.convertSignatureToEIP2098() = %v, want %v", got, tt.want)
			}
		})
	}
}
