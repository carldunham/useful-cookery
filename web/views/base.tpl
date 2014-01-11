<!DOCTYPE HTML PUBLIC "-//IETF//DTD HTML//EN">
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="copyright" content="Copyright &copy; 2014 Carl A. Dunham, All Rights Reserved" />
  <title>{% block title %}{{ title }}{% endblock %}</title>

  <link rel="stylesheet" href="http://code.jquery.com/ui/1.10.3/themes/smoothness/jquery-ui.css">

  <script src="http://ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script>
  <script src="http://ajax.googleapis.com/ajax/libs/jqueryui/1.10.3/jquery-ui.min.js"></script>

  <link rel="stylesheet" href="/css/default.css">
  <script src="/js/cookies.js"></script>
  <script src="/js/main.js"></script>

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
  <span id="loading-indicator"><img src="/images/loading.gif" width="32" height="32"></span>

  <div id="header">
    {% block header %}
    <div id="nav">
      {% block nav %}
      <ul>
	<li><span><a href="/">Home</a></span>
	<li><span><a href="/index">Index</a></span>
	<li><span><a href="/search">Search</a></span>
      </ul>
      {% endblock nav %}
    </div>
    <div class="unitstype">
      Choose Units: <a class="us" href="javascript:chooseUnits('us')">US</a> or <a class="metric" href="javascript:chooseUnits('metric')">Metric</a>
    </div>
    {% endblock header %}
  </div>

  <div id="content">{% block content %}{% endblock content %}</div>

  <div id="footer">
    {% block footer %}
    <div class="copyright">
      <div class="main">Website and original content Copyright &copy; 2014 Carl A. Dunham, All Rights Reserved</div>
      <div class="sub"><i>Usenet Cookbook</i> content from <a href="https://groups.google.com/forum/#!forum/alt.gourmand">alt.gourmand</a> Copyright &copy; 1985-1988 USENET Community Trust</div>
    </div>
    {% endblock footer %}
  </div>
</body>

</html>
