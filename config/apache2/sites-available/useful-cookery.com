
WSGIPythonPath /var/www/useful-cookery.com/current/web/:/var/www/useful-cookery.com/current/lib/

<VirtualHost *:80>
    ServerName www.useful-cookery.com
    ServerAlias uc.carldunham.com
<IfDefine SSL>
    SSLEngine off
</IfDefine>
    DocumentRoot /var/www/useful-cookery.com/current/web/static/
    <Directory "/var/www/useful-cookery.com/current/web/static/">
        AllowOverride All
        Order allow,deny
        Allow from all
    </Directory>
    CustomLog /var/www/useful-cookery.com/log/access.log combined

    WSGIScriptAlias / /var/www/useful-cookery.com/current/web/adapter.wsgi

    <Directory /var/www/useful-cookery.com/current/web/>
        Order allow,deny
        Allow from all
    </Directory>

</VirtualHost>

# funnel everything to the definitive host
#
<VirtualHost *:80>
    ServerName useful-cookery.com
    ServerAlias *.useful-cookery.com
    ServerAlias usefulcookery.com *.usefulcookery.com
    RedirectMatch permanent ^(.*)$ http://www.useful-cookery.com$1
        CustomLog /var/www/useful-cookery.com/log/access-other.log vhost
</VirtualHost>
