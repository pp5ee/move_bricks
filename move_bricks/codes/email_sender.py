import smtplib
from email.mime.text import MIMEText
from email.utils import formataddr


class email_sender():
    def __init__(self):
        self.my_sender = '@qq.com'    # 发件人邮箱账号
        self.my_pass = 'nvmvicy'              # 发件人邮箱密码(当时申请smtp给的口令)
        self.my_user = None
        #self.my_user = '@qq.com'
    def mail(self,message,subject,mail):
        ret=None
        self.my_user=mail
        try:
            print(message,subject)
            msg = MIMEText(str(message), 'plain', 'utf-8')
            # 括号里的对应发件人邮箱昵称、发件人邮箱账号
            msg['From'] = formataddr([str(message), self.my_sender])
            # 括号里的对应收件人邮箱昵称、收件人邮箱账号
            msg['To'] = formataddr(["收件人昵称", self.my_user])
            msg['Subject'] = str(subject)                # 邮件的主题，也可以说是标题

            # 发件人邮箱中的SMTP服务器，端口是465
            server = smtplib.SMTP_SSL("smtp.qq.com", 465)
            server.login(self.my_sender, self.my_pass)  # 括号中对应的是发件人邮箱账号、邮箱密码
            # 括号中对应的是发件人邮箱账号、收件人邮箱账号、发送邮件
            server.sendmail(self.my_sender, [self.my_user, ], msg.as_string())
            server.quit()  # 关闭连接
            ret =True
        except Exception as e:  # 如果 try 中的语句没有执行，则会执行下面的 ret=False
            print(e)
            ret=False
        return ret
