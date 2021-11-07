package importers

import "google.golang.org/api/sheets/v4"

type sheetsImporter struct {
	srv             *sheets.Service // DIP for better testeability. I'm not testing this though
	credentialsPath string
	spreadsheetId   string
}

//func main() {
//ctx := context.Background()
//srv, err := sheets.NewService(ctx, option.WithServiceAccountFile("credentials.json"))
//if err != nil {
//log.Fatalf("Unable to retrieve Sheets client: %v", err)
//}

//// Prints the names and majors of students in a sample spreadsheet:
//// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
//spreadsheetId := "1kZKHDpHl2ZmrkACZ6jrEIAhvC_UyQ_Qjxv2MXNYA-l0"
//readRange := "Expenses!A2:E"
//resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
//if err != nil {
//log.Fatalf("Unable to retrieve data from sheet: %v", err)
//}

//if len(resp.Values) == 0 {
//fmt.Println("No data found.")
//} else {
//fmt.Println("Name, Major:")
//for _, row := range resp.Values {
//// Print columns A and E, which correspond to indices 0 and 4.
//fmt.Printf("%s, %s\n", row[0], row[4])
//}
//}
//}
