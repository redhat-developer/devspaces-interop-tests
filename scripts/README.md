# Run Code Ready Workspaces Test Harness in OSD
To run test harness in OSD first we need to have access to OSD clusters and access to Code Ready Workspaces addon. The
addons are managed in `managed-tenants repo`.

1. Get OFFLINE_TOKEN from OSD (EG: cloud.redhat.com/openshift/token) and put it into `execute-test-harness.sh` file.
2. Get cluster ID from OSD and put it into `osd-test-harness.sh` file
3. Run launch script.

    ```
    ./osd-test-harness.sh
    ```

# Run Code Ready Workspaces Test Harness outside of OSD
1. Access To Openshift Cluster
 

2. Login to the cluster as `admin`

   ```
   oc login -u <user> -p <password> --server=<oc_api_url>
   ```

3. Run the test from your machine

   ```
   ./execute-test-harness.sh <namespace> <report-dir>
   ```

Where are:
 - `namespace` - namespace where you want to deploy test-harness.
 - `report-dir` - Directory where you want to download the results of tests from pods.
