export const RESPONSE_CODE_SUCCESS = 0
export const RESPONSE_CODE_PARAM_ERROR = 1
export const RESPONSE_CODE_DATABASE_ERROR = 2
export const RESPONSE_CODE_NETWORK_ERROR = 3
export const RESPONSE_CODE_SYSTEM_ERROR = 4
export const RESPONSE_CODE_AUTH_ERROR = 5
export const RESPONSE_CODE_NOT_FOUND_ERROR = 6
export const RESPONSE_CODE_UNAUTHORIZED_ERROR = 7
export const RESPONSE_CODE_FILE_ERROR = 8
export const RESPONSE_CODE_FORBIDDEN_ERROR = 9
export const RESPONSE_CODE_AGENT_ERROR = 10
export const RESPONSE_CODE_TOKEN_EXPIRED_ERROR = 11

export const RESPONSE_CODE_MESSAGE_MAP = new Map([
  [RESPONSE_CODE_SUCCESS, 'response_code.success'],
  [RESPONSE_CODE_PARAM_ERROR, 'response_code.param_error'],
  [RESPONSE_CODE_DATABASE_ERROR, 'response_code.database_error'],
  [RESPONSE_CODE_NETWORK_ERROR, 'response_code.network_error'],
  [RESPONSE_CODE_SYSTEM_ERROR, 'response_code.system_error'],
  [RESPONSE_CODE_AUTH_ERROR, 'response_code.auth_error'],
  [RESPONSE_CODE_NOT_FOUND_ERROR, 'response_code.not_found_error'],
  [RESPONSE_CODE_UNAUTHORIZED_ERROR, 'response_code.unauthorized_error'],
  [RESPONSE_CODE_FILE_ERROR, 'response_code.file_error'],
  [RESPONSE_CODE_FORBIDDEN_ERROR, 'response_code.forbidden_error'],
  [RESPONSE_CODE_AGENT_ERROR, 'response_code.agent_error'],
  [RESPONSE_CODE_TOKEN_EXPIRED_ERROR, 'response_code.token_expired_error'],
])

export const RESPONSE_MESSAGE_MAP = new Map([
  ['unauthorized: user not found', 'response_message.user_not_found'],
  ['invalid or expired verification code', 'response_code.verification_code_error'],
  [
    'unauthorized: username or password is incorrect',
    'response_message.username_or_password_is_incorrect',
  ],
])
