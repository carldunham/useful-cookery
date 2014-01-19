{% extends "base.tpl" %}

{% block styles %}
{{ super() }}
  <link rel="stylesheet" href="/css/error404.css">
{% endblock styles %}

{% block content %}
<div class="error404">File not found.</div>

<script type="text/javascript">
  var GOOG_FIXURL_LANG = 'en';
  var GOOG_FIXURL_SITE = 'http://www.useful-cookery.com'
</script>
<script type="text/javascript"
	src="http://linkhelp.clients.google.com/tbproxy/lh/wm/fixurl.js">
</script>
{% endblock %}
