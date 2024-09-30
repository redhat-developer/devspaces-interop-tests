# Run Dev Spaces Test Harness in OpenShift cluster
1. Access To Openshift Cluster
 

2. Login to the cluster as `admin`

   ```
   oc login -u <user> -p <password> --server=<oc_api_url>
   ```

3. Run the test from your machine

   ```
   ./execute-test-harness.sh <user-name>
   ```

Where are:
 - `user-name` - user-name of your OCP cluster.
