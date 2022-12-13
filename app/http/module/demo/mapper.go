package demo

import "github.com/chsir-zy/anan/app/provider/demo"

func UserModelsToUserDTOs(models []UserModel) []UserDTO {
	ret := []UserDTO{}
	for _, model := range models {
		t := UserDTO{
			ID:   model.UserId,
			Name: model.Name,
		}

		ret = append(ret, t)
	}

	return ret
}

func StudentToUserDTOs(students []demo.Student) []UserDTO {
	ret := []UserDTO{}
	for _, student := range students {
		t := UserDTO{
			ID:   student.Id,
			Name: student.Name,
		}

		ret = append(ret, t)
	}

	return ret
}
