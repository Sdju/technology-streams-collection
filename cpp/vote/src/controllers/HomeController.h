//
// Created by zede on 21.05.2023.
//

#ifndef VOTE_SRC_CONTROLLERS_HOMECONTROLLER_H_
#define VOTE_SRC_CONTROLLERS_HOMECONTROLLER_H_

#include "../controllers/BaseController.h"

class HomeController: public BaseController {
 public:
  explicit HomeController(router_type router, ServerData* server_data);

  restinio::request_handling_status_t homeHandler(request_type req, params_type params);
};

#endif //VOTE_SRC_CONTROLLERS_HOMECONTROLLER_H_
