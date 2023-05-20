//
// Created by zede on 20.05.2023.
//

#ifndef VOTE_SRC_SESSION_H_
#define VOTE_SRC_SESSION_H_

#include <string>
#include <memory>
#include "user.h"

struct Session {
  std::string user_agent;
  std::string id;
  std::shared_ptr<User> user;
};

#endif //VOTE_SRC_SESSION_H_
