//
// Created by zede on 20.05.2023.
//

#ifndef VOTE_SRC_BUILDERS_SESSIONBUILDER_H_
#define VOTE_SRC_BUILDERS_SESSIONBUILDER_H_

#include <string>
#include <memory>
#include "../entities/Session.h"

class SessionBuilder {
 public:
  [[nodiscard]] std::shared_ptr<Session> create(const std::string& user_agent, const std::string& id, const std::shared_ptr<User>& user) const;
};

#endif //VOTE_SRC_BUILDERS_SESSIONBUILDER_H_
