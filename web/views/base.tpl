<!DOCTYPE HTML PUBLIC "-//IETF//DTD HTML//EN">
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="copyright" content="Copyright &copy; 2014 Carl A. Dunham, All Rights Reserved" />
  <title>{% block title %}{{ title }}{% endblock %}</title>

  <link rel="stylesheet" href="http://code.jquery.com/ui/1.10.3/themes/smoothness/jquery-ui.css">

  <script src="http://ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script>
  <script src="http://ajax.googleapis.com/ajax/libs/jqueryui/1.10.3/jquery-ui.min.js"></script>

  <link rel="stylesheet" href="css/default.css">
  <script src="js/main.js"></script>

  <style>
    #loading-indicator {
      display: none;
      /*float: right;
      position: absolute;
      top: 0; left: 800px;*/
    }
  </style>

  <script type="text/javascript">

      $(document).ready(function () {

	  $(document).bind("ajaxSend", function() {
	      $("#loading-indicator").show();
	  }).bind("ajaxComplete", function() {
	      $("#loading-indicator").hide();
	  });
	  
      });

  </script>

</head>

<body>
  <span id="loading-indicator"><img src="loading.gif" width="32" height="32"></span>

  <div id="content">{% block content %}{% endblock %}</div>
</body>

</html>