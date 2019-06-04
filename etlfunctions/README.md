# Scenario

A legacy system is delivering batched csv files consisting of customer profiles to a Cloud Storage bucket. This data is un-sanitized and contains potential PII data. We want to clean this data up and import it into our permanent cloud storage bucket (our "data lake") for consumption by other down stream BI/Data pipeline systems. Our default storage format in our data lake is json.

# TODO

- Need to include hash/uuid in final object name to make sure that if +2 batch files arrive at the same time they don't end up with the same final name.
- If processing of the csv fails (i.e. schema change) we should place a temporary object hold on the problematic csv. That would prevent GCE from aging the object out before we get a chance to reprocess it. (https://cloud.google.com/storage/docs/bucket-lock#object-holds)


# Deployment

Create a regional bucket with a 7 day retention policy. Since this bucket contains sensitive uploaded data (not yet GDPR friendly) and serves only as the initial source storage bucket we want to make sure the data is automatically aged out. Given adequate monitoring the 7 day retention should give us enough time to reprocess files should an error occur.

```
gsutil mb -c regional -l us-central1 --retention 7d gs://netlify-churncsv 
```

We also need to create the destination/success bucket - since this is part of our primary "data lake" theres no retention period

```
gsutil mb -c regional -l us-central1 gs://netlify-churncsv-success
```

Create/deploy the function

```
gcloud functions deploy ChurnTransform --runtime go111 --trigger-resource netlify-churncsv --trigger-event google.storage.object.finalize
```

Note: we're not retrying the function on failure - and we're using the default max memory size 

# Improvements

Apache Aarrow's go package has a supposedly nice csv encoder/decoder. It lets you convert records to typed columns (vs just getting a slice of strings from encoding/csv). However, since its relatively new, and im somewhat pressed for time I didn't get a chance to see what the performance is like. 

# Monitoring

# What was difficult for me

I'm not super happy with how I parse/re-encode the csv columns - its not super maintainable. If I have time at the end of the task I'll probably revisit it. 

# Shortcuts Taken

# Cost Estimates

# Notes

Realistically, I don't know if I'd have done this in Go. Python's great at these sort of etl transformation's and clean ups. But Go makes this nice if its likely that this might become a streaming process instead of a batch process - especially if aws-glue/dataflow etc are perhaps unavailable or too costly.