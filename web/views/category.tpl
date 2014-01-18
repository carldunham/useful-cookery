
{% extends "base.tpl" %}

{% block styles %}
{{ super() }}
  <link rel="stylesheet" href="/css/category.css">
{% endblock styles %}

{% block description %}{{ category.name }} {{ super()|lower }}{% endblock description %}
{% block keywords %}{{ category.name|lower }}, {{ super() }}{% endblock keywords %}

{% block content %}
<div class="category">
  {% if error %}
    <div class="error">{{ errorText or 'Category error.' }}</div>
  {% elif recipes %}
    <div class="header">{{ category.name }} Recipes:</div>
    <ol class="recipes">
      {% for recipe in recipes %}
        <li class="recipe" itemscope itemtype="http://schema.org/Recipe">
          <a itemprop="url" href="/recipe/{{ recipe.name }}"><span class="title" itemprop="name">{{ recipe.title }}</span></a> 
	  <span class="description" itemprop="description">{{ recipe.description }}</span>
	</li>
      {% endfor %}
    </ol>
  {% else %}
    <div class="header">No {{ category.name }} recipes found.</div>
  {% endif %}
</div>
{% endblock %}
