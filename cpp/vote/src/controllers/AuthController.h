//
// Created by zede on 21.05.2023.
//

#ifndef VOTE_SRC_CONTROLLERS_AUTHCONTROLLER_H_
#define VOTE_SRC_CONTROLLERS_AUTHCONTROLLER_H_

#include "../controllers/BaseController.h"
#include "../builders/UsersBuilder.h"
#include "../builders/SessionBuilder.h"
#include "../utilies/utilities.h"

class AuthController: public BaseController {
 public:
  explicit AuthController(
      router_type router,
      ServerData* server_data,
      UsersBuilder* users_builder,
      SessionBuilder* session_builder
  );

  restinio::request_handling_status_t authHandler(request_type req, params_type params);

  restinio::request_handling_status_t signupHandler(request_type req, params_type params);

 private:
  UsersBuilder* users_builder_;
  SessionBuilder* session_builder_;
};

#endif //VOTE_SRC_CONTROLLERS_AUTHCONTROLLER_H_
