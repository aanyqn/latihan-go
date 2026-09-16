package model

type Achievement struct {
	ID int `json:"id"`
	StudentID int `json:"student_id"`
	Name string `json:"name"`
	Champion string `json:"champion"`
}

// type CreateAchievementRequest struct {
// 	student_id int `json:"student_id`
// 	name string `json:"name"`
// 	champion string `json:"name"`	
// }

// type ReplaceAchievementRequest struct {
// 	name string `json:"name"`
// 	champion string `json:"name"`	
// }

// type PatchAchievementRequest struct {
// 	name string `json:"name"`
// 	champion string `json:"name"`	
// }