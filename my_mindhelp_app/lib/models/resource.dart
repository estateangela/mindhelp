import 'package:json_annotation/json_annotation.dart';

part 'resource.g.dart';

@JsonSerializable()
class Resource {
  final String id;
  final String name;
  final String type;
  final String address;
  final String phone;
  final String? website;
  final String description;
  final Location? location;
  final List<String> specialties;
  final bool isBookmarked;

  Resource({
    required this.id,
    required this.name,
    required this.type,
    required this.address,
    required this.phone,
    this.website,
    required this.description,
    this.location,
    required this.specialties,
    required this.isBookmarked,
  });

  factory Resource.fromJson(Map<String, dynamic> json) => _$ResourceFromJson(json);
  Map<String, dynamic> toJson() => _$ResourceToJson(this);
}

@JsonSerializable()
class Location {
  final double lat;
  final double lon;

  Location({
    required this.lat,
    required this.lon,
  });

  factory Location.fromJson(Map<String, dynamic> json) => _$LocationFromJson(json);
  Map<String, dynamic> toJson() => _$LocationToJson(this);
}

@JsonSerializable()
class ResourceListResponse {
  final List<Resource> resources;
  final int total;
  final int page;
  final int pageSize;

  ResourceListResponse({
    required this.resources,
    required this.total,
    required this.page,
    required this.pageSize,
  });

  factory ResourceListResponse.fromJson(Map<String, dynamic> json) => 
      _$ResourceListResponseFromJson(json);
  Map<String, dynamic> toJson() => _$ResourceListResponseToJson(this);
}
