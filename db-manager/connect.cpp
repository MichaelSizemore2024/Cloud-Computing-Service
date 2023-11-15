// Connecting to ScyllaDB with a simple C++ program
#include <cassandra.h>
#include <iostream>
#include <memory>

using namespace std;

int main(int argc, char* argv[]) {

  // Allocate the objects that represent cluster and session. 
  CassCluster* cluster = cass_cluster_new();
  CassSession* session = cass_session_new();

  // Add the contact points.
  cass_cluster_set_contact_points(cluster, "scylla");

  // Tries to connect, will block until connect or refused
  auto connect_future = std::unique_ptr<CassFuture, decltype(&cass_future_free)>(cass_session_connect(session, cluster), &cass_future_free);
  if (cass_future_error_code(connect_future.get()) == CASS_OK) {
    cout << "Connected" << endl;
  }
  else{
    //Prints out error code if unable to connect
    cout << "Connection Error: " << cass_error_desc(cass_future_error_code(connect_future.get())) << endl;
    cass_cluster_free(cluster);
    cass_session_free(session);
    return 1;
  }

  // Releases the allocated resources
  cass_cluster_free(cluster);
  cass_session_free(session);
}