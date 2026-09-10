package htmx

import (
	"html/template"
	"jobqueue/entity"
	_interface "jobqueue/interface"
	"net/http"

	"github.com/labstack/echo/v4"
)

// DashboardHandler bridges HTTP (echo) requests to the job service.
// No business logic lives here — only request parsing and response rendering.
type DashboardHandler struct {
	jobService _interface.JobService
}

func NewDashboardHandler(jobService _interface.JobService) *DashboardHandler {
	return &DashboardHandler{jobService: jobService}
}

type variablesForm struct {
	Job1 string
	Job2 string
	Job3 string
}

// Page serves the full HTML page shell.
func (h *DashboardHandler) Page(c echo.Context) error {
	tmpl, err := template.ParseFiles("./web/htmx/dashboard.html")
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return tmpl.Execute(c.Response().Writer, nil)
}

// Message serves the initial fragment: action bar, variables form, and the
// polling containers for status + jobs table.
func (h *DashboardHandler) Message(c echo.Context) error {
	tmpl, err := template.ParseFiles("./web/htmx/message.html")
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	data := variablesForm{Job1: "JobTest1", Job2: "JobTest2", Job3: "JobTest3"}
	return tmpl.Execute(c.Response().Writer, data)
}

// CreateJobs runs SimultaneousCreateJob using Job1/Job2/Job3 from the form.
func (h *DashboardHandler) CreateJobs(c echo.Context) error {
	tasks := []string{
		c.FormValue("job1"),
		c.FormValue("job2"),
		c.FormValue("job3"),
	}
	if _, err := h.jobService.SimultaneousCreateJob(c.Request().Context(), tasks); err != nil {
		return c.String(http.StatusInternalServerError, "failed to create jobs")
	}
	return h.renderJobsTable(c)
}

// UnstableJob runs SimulateUnstableJob (hard-coded task "unstable-job").
func (h *DashboardHandler) UnstableJob(c echo.Context) error {
	if _, err := h.jobService.SimulateUnstableJob(c.Request().Context()); err != nil {
		return c.String(http.StatusInternalServerError, "failed to create unstable job")
	}
	return h.renderJobsTable(c)
}

// Jobs returns the jobs table fragment. Also used by the 2s poll.
func (h *DashboardHandler) Jobs(c echo.Context) error {
	return h.renderJobsTable(c)
}

func (h *DashboardHandler) renderJobsTable(c echo.Context) error {
	jobs, err := h.jobService.GetAllJobs(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to load jobs")
	}
	tmpl, err := template.ParseFiles("./web/htmx/jobs-table.html")
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return tmpl.Execute(c.Response().Writer, struct{ Jobs []*entity.Job }{Jobs: jobs})
}

// Status returns the status summary fragment. Also used by the 2s poll.
func (h *DashboardHandler) Status(c echo.Context) error {
	stats, err := h.jobService.GetAllJobStatus(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to load status")
	}
	tmpl, err := template.ParseFiles("./web/htmx/status-summary.html")
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return tmpl.Execute(c.Response().Writer, stats)
}

// JobDetail returns the detail fragment for one job.
func (h *DashboardHandler) JobDetail(c echo.Context) error {
	id := c.Param("id")
	job, err := h.jobService.GetJobById(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusNotFound, `<div id="job-detail">Job not found</div>`)
	}
	tmpl, err := template.ParseFiles("./web/htmx/job-detail.html")
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return tmpl.Execute(c.Response().Writer, job)
}