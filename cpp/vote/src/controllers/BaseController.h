//
// Created by zede on 20.05.2023.
//

#ifndef VOTE_SRC_CONTROLLERS_BASECONTROLLER_H_
#define VOTE_SRC_CONTROLLERS_BASECONTROLLER_H_

#include "restinio/all.hpp"
#include "../entities/ServerData.h"
#include "../entities/Cookie.h"

class BaseController {
 public:
  using router_type = restinio::router::express_router_t<>*;
  using request_type = std::shared_ptr<restinio::generic_request_t<restinio::no_extra_data_factory_t::data_t>>;
  using params_type = restinio::router::route_params_t;
  using cookies_type = std::unordered_map<std::string, std::string>;

  virtual ~BaseController() = default;

 protected:
  explicit BaseController(router_type router, ServerData* server_data): router_(router), server_data_(server_data) {
  }

  restinio::request_handling_status_t sendResponse(request_type& req, const std::string& answer, const std::string& cookie = "");

  cookies_type* parse_cookies(request_type& req) ;

  std::shared_ptr<Session> load_session(request_type& req, cookies_type* cookies);

  std::string cookie_to_setter(const Cookie& cookie);

  router_type router_;
  ServerData* server_data_;
};

#define addRoute(route, method) this->router_->http_get( \
  route, \
  std::bind(method, this, std::placeholders::_1, std::placeholders::_2))

#endif //VOTE_SRC_CONTROLLERS_BASECONTROLLER_H_
