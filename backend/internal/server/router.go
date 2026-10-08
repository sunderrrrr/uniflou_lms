package server

import (
	"github.com/gin-gonic/gin"
	"github.com/uniflou/uniflou-lms/backend/internal/admin"
	"github.com/uniflou/uniflou-lms/backend/internal/ai"
	"github.com/uniflou/uniflou-lms/backend/internal/auth"
	"github.com/uniflou/uniflou-lms/backend/internal/parent"
	"github.com/uniflou/uniflou-lms/backend/internal/plan"
	"github.com/uniflou/uniflou-lms/backend/internal/progress"
	"github.com/uniflou/uniflou-lms/backend/internal/server/middleware"
	"github.com/uniflou/uniflou-lms/backend/internal/task"
	"github.com/uniflou/uniflou-lms/backend/internal/tutor"
	"github.com/uniflou/uniflou-lms/backend/internal/user"
)

type Handlers struct {
	Auth     *auth.Handler
	User     *user.Handler
	Plan     *plan.Handler
	Task     *task.Handler
	Progress *progress.Handler
	AI       *ai.Handler
	Parent   *parent.Handler
	Tutor    *tutor.Handler
	Admin    *admin.Handler
}

func registerRoutes(r *gin.Engine, h Handlers, jwt *auth.JWTManager, loginLimit, registerLimit int) {
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	v1 := r.Group("/api/v1")
	{
		authG := v1.Group("/auth")
		authG.POST("/register", middleware.RateLimitRegister(registerLimit), h.Auth.Register)
		authG.POST("/login", middleware.RateLimitLogin(loginLimit), h.Auth.Login)
		authG.POST("/refresh", h.Auth.Refresh)
		authG.POST("/reset-password", h.Auth.ResetPassword)
		authG.POST("/reset-password/confirm", h.Auth.ResetPasswordConfirm)

		adminG := v1.Group("")
		adminG.Use(middleware.Auth(jwt), middleware.RequireRoles(auth.RoleAdmin))
		adminG.POST("/admin/tutors", h.Auth.CreateTutor)
		adminG.GET("/admin/users", h.Admin.Users)
		adminG.PATCH("/admin/users/:id/role", h.Admin.UpdateUserRole)
		adminG.DELETE("/admin/users/:id", h.Admin.DeleteUser)
		adminG.GET("/admin/content/subjects", h.Admin.Subjects)
		adminG.POST("/admin/content/subjects", h.Admin.CreateSubject)
		adminG.POST("/admin/content/topics", h.Admin.CreateTopic)
		adminG.GET("/admin/content/tasks", h.Admin.Tasks)
		adminG.POST("/admin/content/tasks", h.Admin.CreateTask)
		adminG.PUT("/admin/content/tasks/:id", h.Admin.UpdateTask)
		adminG.DELETE("/admin/content/tasks/:id", h.Admin.DeleteTask)

		protected := v1.Group("")
		protected.Use(middleware.Auth(jwt))
		{
			protected.POST("/auth/logout", h.Auth.Logout)

			protected.GET("/me", h.User.Me)
			protected.PATCH("/me", h.User.PatchMe)
			protected.POST("/me/link-student", middleware.RequireRoles(auth.RoleParent), h.User.LinkStudent)
			protected.GET("/me/progress", middleware.RequireRoles(auth.RoleStudent), h.Progress.MeProgress)

			student := protected.Group("", middleware.RequireRoles(auth.RoleStudent))
			{
				student.GET("/plan", h.Plan.GetPlan)
				student.PUT("/plan", h.Plan.PutPlan)
				student.GET("/plan/universities", h.Plan.ListUniversities)
				student.POST("/plan/universities", h.Plan.AddUniversity)
				student.DELETE("/plan/universities/:id", h.Plan.DeleteUniversity)
				student.GET("/plan/deadlines", h.Plan.Deadlines)
				student.POST("/plan/deadlines", h.Plan.AddDeadline)
				student.DELETE("/plan/deadlines/:id", h.Plan.DeleteDeadline)

				student.GET("/tasks", h.Task.List)
				student.GET("/tasks/trainer-catalog", h.Task.TrainerCatalog)
				student.GET("/tasks/trainer-topics", h.Task.TrainerTopics)
				student.GET("/tasks/mock-session", h.Task.MockSession)
				student.GET("/tasks/trainer-session", h.Task.TrainingBatch)
				student.GET("/tasks/next", h.Task.Next)
				student.GET("/tasks/next-batch", h.Task.NextBatch)
				student.GET("/tasks/history", h.Task.History)
				student.GET("/tasks/:id", h.Task.Get)
				student.POST("/tasks/:id/attempt", h.Task.Attempt)

				student.GET("/progress", h.Progress.Get)
				student.GET("/progress/subjects", h.Progress.Subjects)
				student.GET("/progress/topics", h.Progress.Topics)
				student.GET("/progress/weak", h.Progress.Weak)
				student.GET("/progress/mock-results", h.Progress.MockResults)
				student.POST("/progress/mock-results", h.Progress.CreateMockResult)

				student.POST("/ai/analyze/task/:attempt_id", h.AI.AnalyzeTask)
				student.POST("/ai/analyze/trainer-session", h.AI.AnalyzeTrainingSession)
				student.POST("/ai/analyze/dashboard", h.AI.AnalyzeDashboard)
				student.POST("/ai/analyze/plan", h.AI.AnalyzePlan)
				student.POST("/ai/analyze/progress", h.AI.AnalyzeProgress)
				student.GET("/ai/reports", h.AI.Reports)
			}

			parentG := protected.Group("/parent", middleware.RequireRoles(auth.RoleParent))
			{
				parentG.GET("/children", h.Parent.Children)
				parentG.GET("/children/:id/progress", h.Parent.ChildProgress)
				parentG.GET("/children/:id/reports", h.Parent.ChildReports)
			}

			tutorG := protected.Group("/tutor", middleware.RequireRoles(auth.RoleTutor, auth.RoleAdmin))
			{
				tutorG.GET("/requests", h.Tutor.Requests)
				tutorG.PATCH("/requests/:id/read", h.Tutor.MarkRequestRead)
				tutorG.GET("/students", h.Tutor.Students)
				tutorG.GET("/students/:id", h.Tutor.StudentDetail)
				tutorG.POST("/students/:id/comments", h.Tutor.Comment)
				tutorG.PUT("/students/:id/target", h.Tutor.SetStudentTarget)
			}

			protected.POST("/tutor/requests", middleware.RequireRoles(auth.RoleStudent, auth.RoleParent), h.Tutor.RequestHelp)
		}
	}
}
