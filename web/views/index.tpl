{% extends "base.tpl" %}

{% block content %}
<div class="rotd">
  <div class=headline>Today's recipe: <span class="name">{{ rotd.name }}</span></div>
  <div class="title">{{ rotd.title }}</div>
  <div class="description">{{ rotd.description }}</div>
</div>
{% endblock %}
