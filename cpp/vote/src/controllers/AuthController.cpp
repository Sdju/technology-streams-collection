//
// Created by zede on 21.05.2023.
//

#include "AuthController.h"

AuthController::AuthController(
    router_type router,
ServerData* server_data,
    UsersBuilder* users_builder,
SessionBuilder* session_builder
)
: BaseController(router, server_data)
, users_builder_(users_builder)
, session_builder_(session_builder) {
  addRoute("/auth/sign-up", &AuthController::signupHandler);
  addRoute("/auth/log-in", &AuthController::authHandler);
}

restinio::request_handling_status_t AuthController::authHandler(request_type req, params_type params) {
  const auto qp = restinio::parse_query(req->header().query());
  auto userName = opt_value<std::string>(qp, "user");
  auto userPassword = opt_value<std::string>(qp, "pwd");

  auto user = server_data_->users_repository_.findUser(userName.value(), userPassword.value());
  if (user == nullptr) {
    fmt::print("AUTH_CONTROLLER::authHandler (failed {})\n", userName);
    return sendResponse(req, "auth failed");
  }

  fmt::print("AUTH_CONTROLLER::authHandler (success {})\n", userName);

  auto userAgent = req->header().get_field_or("User-Agent", "");
  auto session = session_builder_->create(userAgent, generate_uuid(), user);
  server_data_->sessions_repository_.addSession(session);

  Cookie cookie{"SESSION_ID", session->id, true};
  auto cookieSetter = cookie_to_setter(cookie);

  return sendResponse(req, fmt::format("auth success: {}", userName), cookieSetter);
}

restinio::request_handling_status_t AuthController::signupHandler(request_type req, params_type params) {
  const auto qp = restinio::parse_query(req->header().query());
  auto userName = opt_value<std::string>(qp, "user");
  auto userPassword = opt_value<std::string>(qp, "pwd");

  auto user = users_builder_->create(userName.value(), userPassword.value());
  server_data_->users_repository_.addUser(std::move(user));
  server_data_->users_repository_.logTable();

  fmt::print("AUTH_CONTROLLER::signupHandler (user: {})\n", userName);

  return sendResponse(req, fmt::format("name: {} password: {}", userName, userPassword));
}