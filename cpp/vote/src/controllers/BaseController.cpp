//
// Created by zede on 20.05.2023.
//

#include "BaseController.h"

restinio::request_handling_status_t BaseController::sendResponse(request_type& req, const std::string& answer, const std::string& cookie) {
  auto response = req->create_response();
  if (cookie.length() > 0) {
    response.append_header(restinio::http_field::set_cookie, cookie);
  }

  return response
      .append_header(restinio::http_field::content_type, "text/plain; charset=utf-8")
      .set_body(answer)
      .done();
}

BaseController::cookies_type* BaseController::parse_cookies(request_type& req) {
  auto cookiesString = req->header().get_field_or("Cookie", "");
  if (cookiesString.length() == 0) {
    return nullptr;
  }

  auto cookies = new cookies_type();
  std::istringstream f(cookiesString);
  std::string s;
  while (getline(f, s, ';')) {
    std::istringstream f2(s);
    std::string name;
    std::string value;
    getline(f2, name, '=');
    getline(f2, value, '=');
    cookies->emplace(name, value);
  }
  return cookies;
}

std::shared_ptr<Session> BaseController::load_session(request_type& req, cookies_type* cookies) {
  if (cookies == nullptr) {
    return nullptr;
  }
  auto session_iter = cookies->find("SESSION_ID");
  if (session_iter == cookies->end()) {
    return nullptr;
  }
  auto session_id = session_iter->second;
  if (session_id.empty()) {
    return nullptr;
  }

  auto session = server_data_->sessions_repository_.findSession(session_id);
  if (session == nullptr) {
    return nullptr;
  }

  auto userAgent = req->header().get_field_or("User-Agent", "");
  if (session->user_agent != userAgent) {
    return nullptr;
  }

  return session;
}

std::string BaseController::cookie_to_setter(const Cookie& cookie) {
  std::string modifiers = "";

  if (cookie.max_age != -1) {
    modifiers += fmt::format("; Max-Age={}", cookie.max_age);
  }
  if (!cookie.path.empty()) {
    modifiers += fmt::format("; Path={}", cookie.path);
  }
  if (cookie.http_only) {
    modifiers += fmt::format("; HttpOnly");
  }

  return fmt::format("{}={}{}", cookie.name, cookie.value, modifiers);
}
