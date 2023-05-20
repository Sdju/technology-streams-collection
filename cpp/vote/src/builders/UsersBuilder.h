//
// Created by zede on 20.05.2023.
//

#ifndef VOTE_SRC_BUILDERS_USERSBUILDER_H_
#define VOTE_SRC_BUILDERS_USERSBUILDER_H_

#include <string>
#include <memory>
#include "../entities/User.h"

class UsersBuilder {
 public:
  [[nodiscard]] std::shared_ptr<User> create(const std::string& name, const std::string& pwd, int age = -1) const;
};

#endif //VOTE_SRC_BUILDERS_USERSBUILDER_H_
