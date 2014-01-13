{% extends "base.tpl" %}

{% block styles %}
{{ super() }}
  <link rel="stylesheet" href="/css/permindex.css">
{% endblock styles %}

{% block content %}
{% if error or not recipes %}
<div class="error">{{ errorText or 'Recipe not found.' }}</div>
{% else %}
<div class="permindex">
  {% for line in recipes %}
    <div class="line">
      <a href="/recipe/{{ line[3].name }}"><span class="prefix">{{ line[0] }}</span> <span class="key">{{ line[1] }}</span> <span class="postfix">{{ line[2] }}</span></a> 
    </div>
  {% endfor %}
</div>
{% endif %}
{% endblock %}
