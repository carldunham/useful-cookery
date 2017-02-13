{% extends "base.tpl" %}

{% block styles %}
{{ super() }}
  <link rel="stylesheet" href="/css/search.css">
{% endblock styles %}

{% block content %}
<div class="search">
  <div class="searchform">
    <form>
      <label for="query">Search: </label>
      <input type="search" id="query" name="q" size="40" placeholder="Enter Search Term(s)" value="{{ q }}">
      <input type="submit" value="Search">
    </form>
  </div>
  {% if error %}
    <div class="error">{{ errorText or 'Search error.' }}</div>
  {% endif %}
  {% if results %}
    <div class="summary">
      <span class="count">{{ nfound }} result{% if nfound != 1 %}s{% endif %} found{% if nfound > 0 %} (in {{ "%0.3f"|format(timeSeconds) }} seconds):{% endif %}</span>
    </div>
    <ol class="results">
      {% for result in results %}
        <li class="result" itemscope itemtype="http://schema.org/Recipe">
          <a itemprop="url" href="/recipe/{{ result.name }}"><span itemprop="name">{{ result.title }}</span></a>
	        <span class="description" itemprop="description">{{ result.description }}</span>
	        <!-- <span class="score">{{ "%0.2f"|format(result.score) }}</span> -->
	      </li>
      {% endfor %}
    </ol>
  {% endif %}
</div>
{% endblock %}
