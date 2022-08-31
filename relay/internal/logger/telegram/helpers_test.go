package telegram

import (
	"reflect"
	"testing"
)

func Test_safeSplitText(t *testing.T) {
	type args struct {
		text      string
		length    int
		separator string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Simple case",
			args: args{
				text:      "hello world",
				length:    6,
				separator: " ",
			},
			want: []string{"hello", "world"},
		},
		{
			name: "Realistic case",
			args: args{
				text: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Praesent consequat lectus massa, in aliquet lorem lobortis sit amet.
Sed ipsum leo, egestas in rhoncus iaculis, interdum at magna. Sed a magna sit amet nisl vehicula suscipit sed sit amet mi. Morbi varius maximus neque nec pretium.
Nullam urna leo, pretium in lacinia vel, porttitor condimentum tortor. Curabitur condimentum nisi at ullamcorper fermentum.
Praesent malesuada turpis quis varius dictum. Aliquam sed interdum urna.

Ut ultrices orci sed aliquet commodo. Aliquam erat volutpat. Maecenas leo lacus, ornare ac nunc eget, eleifend posuere quam. Integer sodales aliquam egestas. Cras scelerisque, diam a scelerisque interdum, felis nunc scelerisque mauris, ac ornare turpis sapien id velit. Donec mollis, lectus at consectetur consequat, orci lacus ullamcorper quam, et pharetra odio felis a velit. In tristique convallis urna quis dictum. Curabitur condimentum augue sit amet ante bibendum semper. Pellentesque ac nunc a ex ultricies auctor. Ut in vestibulum erat. 
{
	"address": "422 Dekoven Court, Foxworth, Indiana, 8804",
    "about": "Quis cillum occaecat amet exercitation velit aliqua. Tempor cupidatat fugiat ea est est ut aliqua nostrud excepteur laborum ex do excepteur. Dolore laboris velit mollit voluptate ullamco. In nostrud aute dolor sunt labore. Exercitation reprehenderit ipsum qui minim sit. Laborum qui proident consequat ullamco do ea cillum. Reprehenderit ullamco nulla ad elit labore officia minim officia aliquip sint sunt ut occaecat ipsum.\r\n",
}`,
				length:    1000,
				separator: " ",
			},
			want: []string{
				`Lorem ipsum dolor sit amet, consectetur adipiscing elit. Praesent consequat lectus massa, in aliquet lorem lobortis sit amet.
Sed ipsum leo, egestas in rhoncus iaculis, interdum at magna. Sed a magna sit amet nisl vehicula suscipit sed sit amet mi. Morbi varius maximus neque nec pretium.
Nullam urna leo, pretium in lacinia vel, porttitor condimentum tortor. Curabitur condimentum nisi at ullamcorper fermentum.
Praesent malesuada turpis quis varius dictum. Aliquam sed interdum urna.

Ut ultrices orci sed aliquet commodo. Aliquam erat volutpat. Maecenas leo lacus, ornare ac nunc eget, eleifend posuere quam. Integer sodales aliquam egestas. Cras scelerisque, diam a scelerisque interdum, felis nunc scelerisque mauris, ac ornare turpis sapien id velit. Donec mollis, lectus at consectetur consequat, orci lacus ullamcorper quam, et pharetra odio felis a velit. In tristique convallis urna quis dictum. Curabitur condimentum augue sit amet ante bibendum semper. Pellentesque ac nunc a ex`,

				`ultricies auctor. Ut in vestibulum erat. 
{
	"address": "422 Dekoven Court, Foxworth, Indiana, 8804",
    "about": "Quis cillum occaecat amet exercitation velit aliqua. Tempor cupidatat fugiat ea est est ut aliqua nostrud excepteur laborum ex do excepteur. Dolore laboris velit mollit voluptate ullamco. In nostrud aute dolor sunt labore. Exercitation reprehenderit ipsum qui minim sit. Laborum qui proident consequat ullamco do ea cillum. Reprehenderit ullamco nulla ad elit labore officia minim officia aliquip sint sunt ut occaecat ipsum.\r\n",
}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := safeSplitText(tt.args.text, tt.args.length, tt.args.separator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("safeSplitText() = %v, want %v", got, tt.want)
			}
		})
	}
}
