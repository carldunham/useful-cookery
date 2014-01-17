{% extends "base.tpl" %}

{% block styles %}
{{ super() }}
  <link rel="stylesheet" href="/css/index.css">
{% endblock styles %}

{% block content %}
<div class="rotd">
  <div class=headline>Today's recipe: <span class="name"><a href="/recipe/{{ rotd.name }}">{{ rotd.title }}</a></span></div>
  {% if rotd.image %}<div class="image"><img src="{{ rotd.image.url }}" width="{{ rotd.image.width }}" height="{{ rotd.image.height }}" alt="{{ rotd.title }}" /></div>{% endif %}
  <div class="description">{{ rotd.description }}</div>
</div>
{% endblock %}
