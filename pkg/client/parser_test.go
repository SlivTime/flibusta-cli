package client

import (
	"os"
	"path"
	"reflect"
	"testing"
)

func TestParseSearch(t *testing.T) {
	type args struct {
		inputFileName string
	}
	tests := []struct {
		name    string
		args    args
		want    *[]ListItem
		wantErr bool
	}{
		{
			"Index page - no list",
			args{"index.html"},
			nil,
			true,
		},
		{
			"Item - no list",
			args{"item.html"},
			nil,
			true,
		},
		{
			"List",
			args{"list.html"},
			&[]ListItem{
				{
					Title:   "Пелевин и поколение пустоты",
					Authors: []string{"Сергей Полотовский", "Роман Козак"},
					ID:      "1",
				},
				{
					Title:   "\"Нео-пелевин\"",
					Authors: []string{"Вадим Сеновский"},
					ID:      "2",
				},
				{
					Title:   "Виктор Пелевин - Синий фонарь",
					Authors: []string{"Сергей Валерьевич Бережной"},
					ID:      "3",
				},
			},
			false,
		},
		{
			"List with many authors for one book",
			args{"many_authors.html"},
			&[]ListItem{
				{
					Title: "Не только Холмс. Детектив времен Конан Дойла [Антология викторианской детективной новеллы]",
					Authors: []string{
						"Эллен Вуд",
						"Грант Аллен", "Кэтрин Луиза Пиркис", "Израэль Зангвилл", "Артур Моррисон", "Фергюс Хьюм", "Элизабет Томазина Мид-Смит", "Юстас Роберт Бартон", "Мэтью Фиппс Шил", "Роберт Уильям Чамберс", "Мелвилл Дэвиссон Пост", "Матиас Макдоннелл Бодкин", "Гай Ньюэлл Бусби", "Эрнест Уильям Хорнунг",
					},
					ID: "510935",
				},
			},
			false,
		},
		{
			"502",
			args{"502.html"},
			nil,
			true,
		},
		{
			"json",
			args{"empty.json"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullPath := path.Join("testdata/parser", tt.args.inputFileName)
			stream, err := os.Open(fullPath)
			if err != nil {
				t.Errorf("Cannot open test data file: %v", fullPath)
				return
			}
			gotResult, err := ParseSearch(stream)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.want) {
				t.Errorf("ParseSearch() gotResult = %v, want %v", gotResult, tt.want)
			}
		})
	}
}

func TestListItem_String(t *testing.T) {
	type fields struct {
		Title   string
		Authors []string
		ID      string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"No author",
			fields{
				"TestBookTitle",
				[]string{},
				"1",
			},
			"1: TestBookTitle <>",
		},
		{
			"Single author",
			fields{
				"TestBookTitle",
				[]string{"TestAuthor"},
				"1",
			},
			"1: TestBookTitle <TestAuthor>",
		},
		{
			"Multiple authors",
			fields{
				"TestBookTitle",
				[]string{
					"TestAuthor1",
					"TestAuthor2",
					"TestAuthor3",
				},
				"1",
			},
			"1: TestBookTitle <TestAuthor1, TestAuthor2, TestAuthor3>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &ListItem{
				Title:   tt.fields.Title,
				Authors: tt.fields.Authors,
				ID:      tt.fields.ID,
			}
			if got := item.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
