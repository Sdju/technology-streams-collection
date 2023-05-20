#include "restinio/all.hpp"

#include "controllers/HomeController.h"
#include "controllers/AuthController.h"

using namespace restinio;

int main() {
  ServerData server_data;
  UsersBuilder users_builder;
  SessionBuilder session_builder;
  auto router = std::make_unique<router::express_router_t<>>();

  HomeController home_controller(router.get(), &server_data);
  AuthController auth_controller(router.get(), &server_data, &users_builder, &session_builder);

  router->non_matched_request_handler(
      [](auto req) {
        return req->create_response(restinio::status_not_found()).connection_close().done();
      });

  // Launching a server with custom traits.
  struct my_server_traits : public default_single_thread_traits_t {
    using request_handler_t = restinio::router::express_router_t<>;
  };

  fmt::print("Server run on: http://localhost:8080");

  restinio::run(
      restinio::on_this_thread<my_server_traits>()
          .address("localhost")
          .port(8080)
          .request_handler(std::move(router)));

  return 0;
}