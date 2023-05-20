//
// Created by zede on 20.05.2023.
//

#ifndef VOTE_SRC_ENTITIES_SERVERDATA_H_
#define VOTE_SRC_ENTITIES_SERVERDATA_H_

#include "../repositories/UsersRepository.h"
#include "../repositories/SessionsRepository.h"

class ServerData {
 public:
  UsersRepository users_repository_;
  SessionsRepository sessions_repository_;
};

#endif //VOTE_SRC_ENTITIES_SERVERDATA_H_
