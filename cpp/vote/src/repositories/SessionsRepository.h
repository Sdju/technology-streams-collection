//
// Created by zede on 20.05.2023.
//

#ifndef VOTE_SRC_REPOSITORIES_SESSIONSREPOSITORY_H_
#define VOTE_SRC_REPOSITORIES_SESSIONSREPOSITORY_H_

#include "../entities/Session.h"

class SessionsRepository {
 public:
  void addSession(const std::shared_ptr<Session>& session);

  [[nodiscard]] std::shared_ptr<Session> findSession(const std::string& id) const;

 private:
  std::unordered_map<std::string, std::shared_ptr<Session>> sessions_;
};

#endif //VOTE_SRC_REPOSITORIES_SESSIONSREPOSITORY_H_
