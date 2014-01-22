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
        <li class="category"><a class="catlink" href="/category/{{ key }}">{{ value }}</a><a class="random" href="/random/{{ key }}">(random)</a></li>
      {% endfor %}
    </ul>
  {% endif %}
</div>
{% endblock content %}
