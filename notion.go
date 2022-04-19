package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/dstotijn/go-notion"
)

func NotionCreatePage(user User) int {
	var orgs = []notion.SelectOptions{}
	// 複数組織に所属している場合に備えて定義を作る
	for i := 0; i < len(user.Organizations); i++ {
		org := notion.SelectOptions{
			Name: strings.TrimSpace(user.Organizations[i]),
		}
		orgs = append(orgs, org)
	}

	pageId := checkExistsUser(user)
	if pageId != "" {
		updateNotionPage(user, orgs, pageId)
	} else {
		createNotionPage(user, orgs)
	}
	return 1
}

func updateNotionPage(user User, orgs []notion.SelectOptions, pageId string) {
	params := updateUserParams(user, orgs)

	client := notion.NewClient(os.Getenv("NOTION_TOKEN"))
	res, err := client.UpdatePage(context.Background(), pageId, params)
	if err != nil {
		fmt.Println("エラーだよ")
		fmt.Printf("%#v\n", err)
		panic(err)
	}
	fmt.Printf("%#v\n", res)
}

func createNotionPage(user User, orgs []notion.SelectOptions) {
	params := createUserParams(user, orgs)

	client := notion.NewClient(os.Getenv("NOTION_TOKEN"))
	res, err := client.CreatePage(context.Background(), params)
	if err != nil {
		fmt.Println("エラーだよ")
		fmt.Printf("%#v\n", err)
		panic(err)
	}
	fmt.Printf("%#v\n", res)
}

// TOOD: 同姓同名の場合どう判断するか
func checkExistsUser(user User) string {
	client := notion.NewClient(os.Getenv("NOTION_TOKEN"))

	query := &notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			Property: "Name",
			Text: &notion.TextDatabaseQueryFilter{
				Contains: user.Name,
			},
		},
	}
	res, err := client.QueryDatabase(context.Background(), os.Getenv("NOTION_DATABASE_ID"), query)
	if err != nil {
		fmt.Println("QueryDatabase エラーだよ")
		fmt.Printf("%#v\n", err)
		panic(err)
	}
	if len(res.Results) == 0 {

		return ""
	}
	fmt.Printf("%s さんはすでに同名のユーザーが登録されています\n", user.Name)
	// TODO: とりあえず1件目を返す
	return res.Results[0].ID
}

func updateUserParams(user User, orgs []notion.SelectOptions) notion.UpdatePageParams {
	var params notion.UpdatePageParams
	d := notion.DatabasePageProperties{}

	if len(orgs) != 0 {
		d["Organizations"] = notion.DatabasePageProperty{
			MultiSelect: orgs,
		}
	}
	d["Image"] = notion.DatabasePageProperty{
		Files: []notion.File{
			{
				Name: fmt.Sprintf("%s.png", strings.Replace(user.Name, " ", "_", -1)), // 空白を_へ置換
				Type: "external",
				External: &notion.FileExternal{
					URL: user.ImageURL,
				},
			},
		},
	}
	d["Email"] = notion.DatabasePageProperty{
		Email: &user.Email,
	}
	params.DatabasePageProperties = d

	return params
}

func createUserParams(user User, orgs []notion.SelectOptions) notion.CreatePageParams {
	if len(orgs) == 0 {
		return notion.CreatePageParams{
			ParentType: notion.ParentTypeDatabase,
			ParentID:   "45efb68f1a4d4335ac768497c1fc6ff4",
			DatabasePageProperties: &notion.DatabasePageProperties{
				"title": notion.DatabasePageProperty{
					Title: []notion.RichText{
						{
							Text: &notion.Text{
								Content: user.Name,
							},
						},
					},
				},
				"Image": notion.DatabasePageProperty{
					Files: []notion.File{
						{
							Name: fmt.Sprintf("%s.png", strings.Replace(user.Name, " ", "_", -1)), // 空白を_へ置換
							Type: "external",
							External: &notion.FileExternal{
								URL: user.ImageURL,
							},
						},
					},
				},
				"Email": notion.DatabasePageProperty{
					Email: &user.Email,
				},
			},
		}
	} else {
		return notion.CreatePageParams{
			ParentType: notion.ParentTypeDatabase,
			ParentID:   "45efb68f1a4d4335ac768497c1fc6ff4",
			DatabasePageProperties: &notion.DatabasePageProperties{
				"title": notion.DatabasePageProperty{
					Title: []notion.RichText{
						{
							Text: &notion.Text{
								Content: user.Name,
							},
						},
					},
				},
				"Organizations": notion.DatabasePageProperty{
					MultiSelect: orgs,
				},
				"Image": notion.DatabasePageProperty{
					Files: []notion.File{
						{
							Name: fmt.Sprintf("%s.png", strings.Replace(user.Name, " ", "_", -1)), // 空白を_へ置換
							Type: "external",
							External: &notion.FileExternal{
								URL: user.ImageURL,
							},
						},
					},
				},
				"Email": notion.DatabasePageProperty{
					Email: &user.Email,
				},
			},
		}
	}
}
