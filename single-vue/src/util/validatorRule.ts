// 表单校验规则
export const loginRule = {
    username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
    ],
        password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
    ],
        confirmPassword: [
        { required: true, message: '请确认密码', trigger: 'blur' },
    ],
        email: [
        { required: true, message: '请输入邮箱', trigger: 'blur' },
        {
            type: 'email',
            message: '请输入正确的邮箱地址',
            trigger: ['blur']
        }
    ],
        code: [
        { required: true, message: '请输入验证码', trigger: 'blur' },
        {
            pattern: /^\d{4,6}$/,
            message: '验证码为4-6位数字',
            trigger: 'blur'
        }
    ],
}

