package splitwise

import (
	"context"
	"encoding/json"
	"fmt"
	_ "log"
	"net/http"
	"time"
)

// Expenses contains method to work with expense resource
type Expenses interface {
	// ExpensesByDate returns current user's expenses
	ExpensesByDate(ctx context.Context, dated_after string, dated_before string) (*Exps, error)
	ExpensesByDateGroup(ctx context.Context, dated_after string, dated_before string, group_name string) (*Exps, error)

	// ExpenseByID returns information of an expense identified by id argument
	// ExpenseByID(ctx context.Context, id uint64) (*Expense, error)

	// // CreateExpense Creates an expense. You may either split an expense equally (only with group_id provided), or
	// // supply a list of shares.
	// //If providing a list of shares, each share must include paid_share and owed_share, and must be identified by one
	// // of the following:
	// //email, first_name, and last_name
	// //user_id
	// //Note: 200 OK does not indicate a successful response. The operation was successful only if errors is empty.
	// CreateExpense(ctx context.Context, dto *CreateCommentDTO) ([]Expense, error)
}

type Expense struct {
	Cost        string    `json:"cost"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Category    struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	Users      []ExpenseUser `json:"users"`
	Details    string        `json:"details"`
	Comments   []Comment     `json:"comments"`
	Deleted_at time.Time     `json:"deleted_at"`
}

type ExpenseUser struct {
	User_id    uint64 `json:"user_id"`
	Owed_share string `json:"owed_share"`
}

type CreateExpenseDTO struct {
}

type Exps struct {
	Exps []Expense `json:"expenses"`
}

func (c SClient) ExpensesByDate(ctx context.Context, dated_after string, dated_before string) (*Exps, error) {
	url := fmt.Sprintf("%s/api/v3.0/get_expenses?limit=0&dated_after=%s&dated_before=%s", c.BaseURL, dated_after, dated_before)
	return c.get_expenses(ctx, url)
}

func (c SClient) ExpensesByDateGroup(ctx context.Context, dated_after string, dated_before string, group_name string) (*Exps, error) {
	g, err := c.Groups(ctx)
	if err != nil {
		return nil, err
	}

	var group_id int
	for _, v := range g {
		if v.Name == group_name {
			group_id = int(v.ID)
		}
	}

	url := fmt.Sprintf("%s/api/v3.0/get_expenses?limit=0&dated_after=%s&dated_before=%s&group_id=%d", c.BaseURL, dated_after, dated_before, group_id)
	return c.get_expenses(ctx, url)
}

func (c SClient) get_expenses(ctx context.Context, url string) (*Exps, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	var response Exps
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
