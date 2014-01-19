{% extends "base.tpl" %}

{% block styles %}
{{ super() }}
  <link rel="stylesheet" href="/css/index.css">
{% endblock styles %}

{% block content %}
<div class="rotd" itemscope itemtype="http://schema.org/Recipe">
  <div class=headline>Featured Recipe: <span class="name"><a itemprop="url" href="/recipe/{{ rotd.name }}"><span itemprop="name">{{ rotd.title }}</span></a></span></div>
  {% if rotd.image %}<div class="image"><img itemprop="image" src="{{ rotd.image.url }}" width="{{ rotd.image.width|int }}" height="{{ rotd.image.height|int }}" alt="{{ rotd.title }}" /></div>{% endif %}
  <div class="description" itemprop="description">{{ rotd.description }}</div>
</div>
{% endblock %}
