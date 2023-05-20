//
// Created by zede on 21.05.2023.
//

#include "HomeController.h"

HomeController::HomeController(router_type router, ServerData* server_data): BaseController(router, server_data) {
  addRoute("/", &HomeController::homeHandler);
}

restinio::request_handling_status_t HomeController::homeHandler(request_type req, params_type params) {
  auto cookies = parse_cookies(req);
  auto session = load_session(req, cookies);

  fmt::print("HOME_CONTROLLER::homeHandler\n");

  if (session != nullptr) {
    return sendResponse(req, fmt::format("Hello, {}!", session->user->name));
  }

  return sendResponse(req, "Login to see greetings!");
}