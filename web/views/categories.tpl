{% extends "base.tpl" %}

{% block styles %}
{{ super() }}
  <link rel="stylesheet" href="/css/categories.css">
{% endblock styles %}

{% block content %}
  {% if error %}
    <div class="error">{{ errorText or 'Unknown error.' }}</div>
  {% else %}
    <ul class="categories">
      {% for key,value in categories|dictsort(false, 'value') %}
        <li class="category"><a href="/category/{{ key }}">{{ value }}</a></li>
      {% endfor %}
    </ul>
  {% endif %}
</div>
{% endblock content %}
