package service

import (
	"context"
	"server/internal/repository"
	"time"
)

func StartCleanupTask(repo *repository.AuthRepository){
	go func(){
		for{
			ctx := context.Background()
			repo.DB.Exec(ctx,"delete from refresh_tokens where expires_at < NOW()")
			repo.DB.Exec(ctx,"delete from jwt_denylist where expires_at < NOW()")
			time.Sleep(1*time.Hour)
		}
	}()
}