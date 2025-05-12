# EWA API 文档

## 基础信息
- 基础URL: `http://localhost:8080`
- 所有需要认证的接口都需要在请求头中携带 `Authorization: Bearer <token>`
- 响应格式统一为 JSON

## 认证相关接口

### 用户注册
- **URL**: `/api/register`
- **方法**: `POST`
- **描述**: 注册新用户
- **请求体**:
```json
{
    "username": "string",    // 用户名，3-50个字符
    "password": "string",    // 密码，6-50个字符
    "email": "string",       // 邮箱地址
    "phone": "string",       // 手机号，11位
    "nickname": "string"     // 昵称，可选
}
```
- **响应**:
```json
{
    "message": "User registered successfully"
}
```

### 用户登录
- **URL**: `/api/login`
- **方法**: `POST`
- **描述**: 用户登录并获取访问令牌
- **请求体**:
```json
{
    "username": "string",    // 用户名
    "password": "string"     // 密码
}
```
- **响应**:
```json
{
    "token": "string",       // JWT token
    "user": {
        "id": "number",
        "username": "string",
        "email": "string",
        "phone": "string",
        "nickname": "string",
        "avatar": "string"
    }
}
```

## 用户相关接口

### 更新用户信息
- **URL**: `/api/profile`
- **方法**: `PUT`
- **描述**: 更新当前登录用户的信息
- **认证**: 需要
- **请求体**:
```json
{
    "nickname": "string",    // 昵称，可选
    "avatar": "string",      // 头像URL，可选
    "phone": "string",       // 手机号，可选，11位
    "email": "string"        // 邮箱，可选
}
```
- **响应**:
```json
{
    "message": "Profile updated successfully",
    "user": {
        "id": "number",
        "username": "string",
        "email": "string",
        "phone": "string",
        "nickname": "string",
        "avatar": "string"
    }
}
```

## 模型相关接口

### 创建模型
- **URL**: `/api/models`
- **方法**: `POST`
- **描述**: 创建新的模型
- **认证**: 需要
- **请求体**:
```json
{
    "name": "string",        // 模型名称，最大100字符
    "version": "string",     // 版本号，最大50字符
    "file_url": "string",    // 模型文件URL
    "file_size": "number",   // 文件大小
    "file_type": "string",   // 文件类型，最大20字符
    "description": "string", // 描述，可选
    "category": "string",    // 分类，最大50字符
    "tags": "string"         // 标签，最大255字符，可选
}
```
- **响应**:
```json
{
    "message": "Model created successfully",
    "model": {
        "id": "number",
        "name": "string",
        "version": "string",
        "file_url": "string",
        "file_size": "number",
        "file_type": "string",
        "description": "string",
        "category": "string",
        "tags": "string",
        "status": "number",
        "download_count": "number",
        "created_by": "number",
        "created_at": "string",
        "updated_at": "string"
    }
}
```

### 更新模型
- **URL**: `/api/models/:id`
- **方法**: `PUT`
- **描述**: 更新模型信息
- **认证**: 需要（需要编辑权限）
- **请求体**:
```json
{
    "name": "string",        // 模型名称，可选，最大100字符
    "version": "string",     // 版本号，可选，最大50字符
    "description": "string", // 描述，可选
    "category": "string",    // 分类，可选，最大50字符
    "tags": "string",        // 标签，可选，最大255字符
    "status": "number"       // 状态，可选，0或1
}
```
- **响应**:
```json
{
    "message": "Model updated successfully",
    "model": {
        // 同创建模型的响应
    }
}
```

### 删除模型
- **URL**: `/api/models/:id`
- **方法**: `DELETE`
- **描述**: 删除模型
- **认证**: 需要（需要超级权限）
- **响应**:
```json
{
    "message": "Model deleted successfully"
}
```

### 获取模型详情
- **URL**: `/api/models/:id`
- **方法**: `GET`
- **描述**: 获取模型详细信息
- **认证**: 需要（需要查看权限）
- **响应**:
```json
{
    "model": {
        // 同创建模型的响应
    }
}
```

### 获取模型列表
- **URL**: `/api/models`
- **方法**: `GET`
- **描述**: 获取模型列表
- **认证**: 需要
- **查询参数**:
  - `page`: 页码，默认1
  - `page_size`: 每页数量，默认10
  - `category`: 分类筛选，可选
  - `keyword`: 关键词搜索，可选
- **响应**:
```json
{
    "models": [
        {
            // 同创建模型的响应
        }
    ],
    "total": "number"
}
```

### 获取模型文件
- **URL**: `/api/models/:id/file`
- **方法**: `GET`
- **描述**: 获取模型文件内容
- **认证**: 需要（需要查看权限）
- **响应**:
```json
{
    "file_url": "string",    // 模型文件的临时访问URL
    "file_name": "string",   // 文件名
    "file_size": "number",   // 文件大小（字节）
    "file_type": "string",   // 文件类型
    "expire_time": "string"  // URL过期时间
}
```

### 获取模型文件上传URL
- **URL**: `/api/models/upload-url`
- **方法**: `POST`
- **描述**: 获取模型文件上传的临时URL
- **认证**: 需要
- **请求体**:
```json
{
    "file_name": "string",   // 文件名
    "file_type": "string"    // 文件类型
}
```
- **响应**:
```json
{
    "upload_url": "string",  // 文件上传的临时URL
    "file_url": "string",    // 上传后的文件访问URL
    "expire_time": "string"  // URL过期时间
}
```

## 权限相关接口

### 授予权限
- **URL**: `/api/permissions`
- **方法**: `POST`
- **描述**: 为用户授予模型权限
- **认证**: 需要（需要管理权限）
- **请求体**:
```json
{
    "user_id": "number",           // 用户ID
    "model_id": "number",          // 模型ID
    "permission_level": "number",  // 权限级别(1-4)
    "is_owner": "number",          // 是否所有者，可选，0或1
    "expire_time": "string",       // 过期时间，可选
    "grant_reason": "string"       // 授权原因，可选，最大255字符
}
```
- **响应**:
```json
{
    "message": "Permission granted successfully",
    "permission": {
        "id": "number",
        "user_id": "number",
        "model_id": "number",
        "permission_level": "number",
        "is_owner": "number",
        "expire_time": "string",
        "granted_by": "number",
        "grant_reason": "string",
        "status": "number",
        "created_at": "string",
        "updated_at": "string"
    }
}
```

### 更新权限
- **URL**: `/api/permissions/:id`
- **方法**: `PUT`
- **描述**: 更新权限信息
- **认证**: 需要（需要管理权限）
- **请求体**:
```json
{
    "permission_level": "number",  // 权限级别，可选，1-4
    "is_owner": "number",          // 是否所有者，可选，0或1
    "expire_time": "string",       // 过期时间，可选
    "status": "number"             // 状态，可选，0或1
}
```
- **响应**:
```json
{
    "message": "Permission updated successfully",
    "permission": {
        // 同授予权限的响应
    }
}
```

### 撤销权限
- **URL**: `/api/permissions/:id`
- **方法**: `DELETE`
- **描述**: 撤销用户对模型的权限
- **认证**: 需要（需要管理权限）
- **响应**:
```json
{
    "message": "Permission revoked successfully"
}
```

### 获取权限列表
- **URL**: `/api/models/:model_id/permissions`
- **方法**: `GET`
- **描述**: 获取模型的所有权限记录
- **认证**: 需要（需要管理权限）
- **响应**:
```json
{
    "permissions": [
        {
            // 同授予权限的响应
        }
    ]
}
```

## 权限级别说明
1. 查看权限：可以查看模型信息
2. 编辑权限：可以更新模型信息
3. 管理权限：可以管理模型权限
4. 超级权限：可以删除模型

## 错误响应
所有接口在发生错误时都会返回以下格式的响应：
```json
{
    "error": "错误信息"
}
```

常见错误状态码：
- 400: 请求参数错误
- 401: 未认证
- 403: 权限不足
- 404: 资源不存在
- 500: 服务器内部错误 