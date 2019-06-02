# netlify (WIP)

# components

* churnprofiles - package to work with the churn profile's and files
* etlfunctions - the lone etl function that gets invoked as a cloud function and converts a csv to a cleaned json file
* services - the single worker that enriches the contained churn profiles with some churn risk properties and stores them for querying. 