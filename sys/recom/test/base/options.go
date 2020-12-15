package base

type RuntimeOptions struct {
	Verbose bool // Verbose switch
	FitJobs int  // Number of jobs for model fitting
	CVJobs  int  // Number of jobs for cross validation
}
