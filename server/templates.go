package server

const userTemplate = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>Пользователь</title>
</head>
<body>

    <div id="user">
        <p id="login">{{.Login}}</p>
        <p id="first_name">{{.FirstName}}</p>
        <p id="last_name">{{.LastName}}</p>
        <p id="age">{{.Age}}</p>
        <p id="sex">{{.Sex}}</p>
        <p id="city">{{.City}}</p>
        <p id="hobby">{{.Hobby}}</p>
        <p id="friends">{{range .Friends}}<div>{{ . }}</div>{{end}}</p>
        <input type="button" value="Подружиться" id="make_friends"/>
    </div>

	<script>
		var nikname = {{.Login}};
		document.getElementById("make_friends").addEventListener("click", function () {
			fetch('/api/make_friends', {
				headers: { "Content-Type": "application/json; charset=utf-8"},
				method: 'POST',
				credentials: 'include',
				body: JSON.stringify({friend: nikname})
			});
		});
	</script>
</body>`

const loginTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Login</title>
</head>
<body>
    <p>Вход или регистрация</p>
    <form action="/login">
        <p><input name="login"> <input type="password" name="password"></p>
        <p><input type="submit" value="Войти"></p>
    </form>
    <p><a  href="registration">Регистрация</a></p>

    <script>
 		document.getElementsByTagName("form")[0].addEventListener("submit", function (e) {
            const requestData = toJSONString(this);
            fetch('/api/login', {
                headers: { "Content-Type": "application/json; charset=utf-8" },
                method: 'POST',
                body: requestData
            })
                .then(response => response.json())
                .then(data => {
					const login = JSON.parse(requestData).login;

                    createCookie("token", data.token, 1);
                    createCookie("login", login, 1);
                    window.location = "/users/"+login;
                });
            e.preventDefault();
        });

        function toJSONString( form ) {
            var obj = {};
            var elements = form.querySelectorAll( "input, select, textarea" );
            for( var i = 0; i < elements.length; ++i ) {
                var element = elements[i];
                var name = element.name;
                var value = element.value;

                if( name ) {
                    obj[ name ] = value;
                }
            }

            return JSON.stringify( obj );
        }

        function createCookie(name, value, days) {
            var expires;
            if (days) {
                var date = new Date();
                date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
                expires = "; expires=" + date.toGMTString();
            }
            else {
                expires = "";
            }
            document.cookie = name + "=" + value + expires + "; path=/";
        }
        function getCookie(c_name) {
            if (document.cookie.length > 0) {
                c_start = document.cookie.indexOf(c_name + "=");
                if (c_start != -1) {
                    c_start = c_start + c_name.length + 1;
                    c_end = document.cookie.indexOf(";", c_start);
                    if (c_end == -1) {
                        c_end = document.cookie.length;
                    }
                    return unescape(document.cookie.substring(c_start, c_end));
                }
            }
            return "";
        }
    </script>
</body>
`

const registrationTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Регистрация</title>
</head>
<body>
    <form action="/new">
        <p>
            <input name="login" placeholder="Логин">
            <input name="password" placeholder="Пароль" type="password" >
            <input name="first_name" placeholder="Имя">
            <input name="last_name" placeholder="Фамилия">
            <input name="age" placeholder="Возраст">
            <input name="sex" placeholder="Пол">
            <input name="city" placeholder="Город">
            <input name="hobby" placeholder="Хобби">
        </p>
        <p><input type="submit" value="Зарегистрироваться"></p>
    </form>

    <script>
        document.getElementsByTagName("form")[0].addEventListener("submit", function (e) {
            const requestData = toJSONString(this);
            fetch('/api/new', {
                headers: { "Content-Type": "application/json; charset=utf-8" },
                method: 'POST',
                body: requestData
            })
                .then(response => response.json())
                .then(data => {
					const login = JSON.parse(requestData).login;

                    createCookie("token", data.token, 1);
                    createCookie("login", login, 1);
                    window.location = "/users/"+login;
                });
            e.preventDefault();
        });

        function toJSONString( form ) {
            var obj = {};
            var elements = form.querySelectorAll( "input, select, textarea" );
            for( var i = 0; i < elements.length; ++i ) {
                var element = elements[i];
                var name = element.name;
                var value = element.value;

                if( name ) {
                    obj[ name ] = value;
                }
            }

            return JSON.stringify( obj );
        }

        function createCookie(name, value, days) {
            var expires;
            if (days) {
                var date = new Date();
                date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
                expires = "; expires=" + date.toGMTString();
            }
            else {
                expires = "";
            }
            document.cookie = name + "=" + value + expires + "; path=/";
        }

        function getCookie(c_name) {
            if (document.cookie.length > 0) {
                c_start = document.cookie.indexOf(c_name + "=");
                if (c_start != -1) {
                    c_start = c_start + c_name.length + 1;
                    c_end = document.cookie.indexOf(";", c_start);
                    if (c_end == -1) {
                        c_end = document.cookie.length;
                    }
                    return unescape(document.cookie.substring(c_start, c_end));
                }
            }
            return "";
        }
    
    </script>
</body>
</html>`
