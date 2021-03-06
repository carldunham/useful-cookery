<!DOCTYPE HTML PUBLIC "-//IETF//DTD HTML//EN">
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="copyright" content="Copyright &copy; 2014-2017 Carl A. Dunham, All Rights Reserved" />
  <title>{% block title %}{{ title }}{% endblock %}</title>

  <meta name="description" content="{% block description %}Recipes from the original USENET Cookbook (and more).{% endblock %}" />
  <meta name="keywords" content="{% block keywords %}recipes, food, cooking, delicious, cooking, dinner, lunch, vegetarian, usenet, cookbook{% endblock %}" />

  {% block styles %}
  <link rel="stylesheet" href="/css/base.css">
  {% endblock styles %}

  {% block javascript %}
  <script src="/js/cookies.js"></script>
  <script src="/js/main.js"></script>
  {% endblock javascript %}

</head>

<body>
  <div id="header">
    {% block header %}
    <div class="bar"></div>
    <div class="nav">
      {% block nav %}
      <ul>
	<li><span><a href="/">Home</a></span></li>
	<li><span><a href="/categories">Categories</a></span></li>
	<li><span><a href="/index">Index</a></span></li>
	<li><span><a href="/search">Search</a></span></li>
	<li><span><a href="/random">Random</a></span></li>
      </ul>
      {% endblock nav %}
    </div>
    {# <div class="logo"><span class="name">Useful Cookery</span><span class="tagline">Recipes restored for the global village</span></div> #}
    <div class="logo"><img src="/images/logo-med.png" alt="Useful Cookery - Recipes restored for the global village" /></div>
    <div class="unitstype">Choose Units: <a class="us" href="javascript:chooseUnits('us')">US</a> or <a class="metric" href="javascript:chooseUnits('metric')">Metric</a></div>
    {% endblock header %}
  </div>

  <div class="clear"></div>

  <div id="content">{% block content %}{% endblock content %}</div>

  <div class="clear"></div>

  <div id="footer">
    {% block footer %}
    <div class="copyright">
      <div class="main">Website and original content Copyright &copy; 2014-2017 Carl A. Dunham, All Rights Reserved</div>
      <div class="sub"><i>Usenet Cookbook</i> content from <a href="https://groups.google.com/forum/#!forum/alt.gourmand">alt.gourmand</a> Copyright &copy; 1985-1988 USENET Community Trust</div>
    </div>
    {% endblock footer %}
  </div>

  <script>
    (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
    (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
    m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
    })(window,document,'script','//www.google-analytics.com/analytics.js','ga');

    ga('create', 'UA-47126324-1', 'useful-cookery.com');
    ga('require', 'linkid', 'linkid.js');
    {% block tracking_send %}ga('send', 'pageview');{% endblock tracking_send %}
  </script>

</body>

</html>
