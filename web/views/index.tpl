{% extends "base.tpl" %}

{% block content %}
<div class="rotd">
  <div class=headline>Today's recipe: <span class="name"><a href="/recipe/{{ rotd.name }}">{{ rotd.name }}</a></span></div>
  <div class="title">{{ rotd.title }}</div>
  <div class="description">{{ rotd.description }}</div>
</div>
{% endblock %}
