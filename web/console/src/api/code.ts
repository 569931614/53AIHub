// 使用枚举替代常量
export const ResponseCode = {
  SUCCESS: 0,
  PARAM_ERROR: 1,
  DATABASE_ERROR: 2,
  NETWORK_ERROR: 3,
  SYSTEM_ERROR: 4,
  AUTH_ERROR: 5,
  NOT_FOUND_ERROR: 6,
  UNAUTHORIZED_ERROR: 7,
  FILE_ERROR: 8,
  FORBIDDEN_ERROR: 9,
  AGENT_ERROR: 10,
  TOKEN_EXPIRED_ERROR: 11,
  VERIFICATION_CODE_ERROR: -14,
}

export const ResponseStatus = {
  SUCCESS: 200,
  CREATED: 201,
  NO_CONTENT: 204,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  INTERNAL_SERVER_ERROR: 500,
  BAD_GATEWAY: 502,
  SERVICE_UNAVAILABLE: 503,
  GATEWAY_TIMEOUT: 504,
}

export const RESPONSE_CODE_MESSAGE_MAP = new Map([
  [ResponseCode.SUCCESS, 'response_code.success'],
  [ResponseCode.PARAM_ERROR, 'response_code.param_error'],
  [ResponseCode.DATABASE_ERROR, 'response_code.database_error'],
  [ResponseCode.NETWORK_ERROR, 'response_code.network_error'],
  [ResponseCode.SYSTEM_ERROR, 'response_code.system_error'],
  [ResponseCode.AUTH_ERROR, 'response_code.auth_error'],
  [ResponseCode.NOT_FOUND_ERROR, 'response_code.not_found_error'],
  [ResponseCode.UNAUTHORIZED_ERROR, 'response_code.unauthorized_error'],
  [ResponseCode.FILE_ERROR, 'response_code.file_error'],
  [ResponseCode.FORBIDDEN_ERROR, 'response_code.forbidden_error'],
  [ResponseCode.AGENT_ERROR, 'response_code.agent_error'],
  [ResponseCode.TOKEN_EXPIRED_ERROR, 'response_code.token_expired_error'],
  [ResponseCode.VERIFICATION_CODE_ERROR, 'response_code.verification_code_error'],
])

// 错误码映射
export const ERROR_MESSAGES = new Map([
  [ResponseStatus.BAD_REQUEST, 'response_status.400'],
  [ResponseStatus.UNAUTHORIZED, 'response_status.401'],
  [ResponseStatus.FORBIDDEN, 'response_status.403'],
  [ResponseStatus.NOT_FOUND, 'response_status.404'],
  [ResponseStatus.SERVICE_UNAVAILABLE, 'response_status.500'],
  [ResponseStatus.BAD_GATEWAY, 'response_status.502'],
  [ResponseStatus.INTERNAL_SERVER_ERROR, 'response_status.503'],
  [ResponseStatus.GATEWAY_TIMEOUT, 'response_status.504'],
])

export const RESPONSE_MESSAGE_MAP = new Map([
  ['unauthorized: user not found', 'response_message.user_not_found'],
  [
    'unauthorized: username or password is incorrect',
    'response_message.username_or_password_is_incorrect',
  ],
  ['auth failed: record not found', 'response_message.user_not_found'],
  ['invalid or expired verification code', 'response_code.verification_code_error'],
  ['operate too fast', 'response_message.operate_too_fast'],
])

export const RESPONSE_DATA_MESSAGE_MAP = new Map([
  ['该域名已被绑定', 'response_data.domain_already_bound'],
  ['您已有一个正在审核中的申请', 'response_data.apply_already_submitted'],
])
